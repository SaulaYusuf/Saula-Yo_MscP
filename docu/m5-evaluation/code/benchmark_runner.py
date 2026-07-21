"""
M5: Throughput & Latency Benchmark
Runs the sensor ingestion at different concurrency levels and records performance.
"""

import asyncio
import csv
import time
import logging
import json
import statistics
from pathlib import Path

import aiohttp

# Use the same sensor ingestion logic, but with configurable concurrency
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

SENSOR_URL = "http://localhost:8080/api/sensor"
DATASET_PATH = Path(__file__).parents[3] / "data" / "raw" / "shipment-sensor-dataset.csv"

async def send_sensor(session, row, start_times, commit_times):
    payload = {
        "sensor_id": row.get("Sensor_ID", "unknown"),
        "temp_c": float(row.get("Temp_Max_C", 0.0)),
        "humidity": float(row.get("Humidity_Pct", 0.0)),
        "timestamp": row.get("Timestamp", ""),
    }
    start = time.time()
    try:
        async with session.post(SENSOR_URL, json=payload, timeout=10) as resp:
            end = time.time()
            start_times.append(start)
            commit_times.append(end)  # approximate (we get HTTP response)
            if resp.status == 200:
                return True
            else:
                return False
    except Exception:
        return False

async def run_benchmark(concurrency, total_records=8000):
    rows = []
    with open(DATASET_PATH, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            rows.append(row)
    rows = rows[:total_records]  # limit

    queue = asyncio.Queue()
    for row in rows:
        await queue.put(row)
    for _ in range(concurrency):
        await queue.put(None)

    start_times = []
    commit_times = []
    results = []

    async def worker(session):
        while True:
            row = await queue.get()
            if row is None:
                queue.task_done()
                break
            ok = await send_sensor(session, row, start_times, commit_times)
            results.append(ok)
            queue.task_done()

    async with aiohttp.ClientSession() as session:
        workers = [asyncio.create_task(worker(session)) for _ in range(concurrency)]
        await queue.join()
        for w in workers:
            w.cancel()

    success = sum(1 for r in results if r)
    total_time = commit_times[-1] - start_times[0] if commit_times and start_times else 0
    tps = len(rows) / total_time if total_time > 0 else 0
    latencies = [commit_times[i] - start_times[i] for i in range(len(start_times)) if i < len(commit_times)]
    if latencies:
        p50 = statistics.median(latencies)
        p95 = statistics.quantiles(latencies, n=20)[18] if len(latencies) >= 20 else max(latencies)
        p99 = statistics.quantiles(latencies, n=100)[98] if len(latencies) >= 100 else max(latencies)
    else:
        p50 = p95 = p99 = 0

    return {
        "concurrency": concurrency,
        "records": len(rows),
        "total_time": total_time,
        "tps": tps,
        "success_rate": success / len(rows) * 100,
        "latency_p50": p50,
        "latency_p95": p95,
        "latency_p99": p99,
    }

async def main():
    concurrency_levels = [10, 50, 100, 200]
    results = []
    for c in concurrency_levels:
        logger.info(f"Running benchmark with concurrency={c}")
        res = await run_benchmark(c)
        results.append(res)
        logger.info(f"Concurrency {c}: TPS={res['tps']:.2f}, p50={res['latency_p50']:.3f}s, p95={res['latency_p95']:.3f}s")

    # Save results to JSON
    with open("benchmark_results.json", "w") as f:
        json.dump(results, f, indent=2)
    logger.info("Results saved to benchmark_results.json")

if __name__ == "__main__":
    asyncio.run(main())