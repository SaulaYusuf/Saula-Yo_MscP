import json
import matplotlib.pyplot as plt

with open("benchmark_results.json", "r") as f:
    data = json.load(f)

concurrency = [d["concurrency"] for d in data]
tps = [d["tps"] for d in data]
p50 = [d["latency_p50"] for d in data]
p95 = [d["latency_p95"] for d in data]
p99 = [d["latency_p99"] for d in data]

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(12, 5))

# TPS vs Concurrency
ax1.plot(concurrency, tps, marker='o', color='blue')
ax1.set_xlabel('Concurrency (workers)')
ax1.set_ylabel('Throughput (TPS)')
ax1.set_title('Throughput vs. Concurrency')
ax1.grid(True)

# Latency Percentiles
ax2.plot(concurrency, p50, marker='s', label='p50', color='green')
ax2.plot(concurrency, p95, marker='^', label='p95', color='orange')
ax2.plot(concurrency, p99, marker='d', label='p99', color='red')
ax2.set_xlabel('Concurrency (workers)')
ax2.set_ylabel('Latency (seconds)')
ax2.set_title('Latency Percentiles vs. Concurrency')
ax2.legend()
ax2.grid(True)

plt.tight_layout()
plt.savefig('benchmark_plots.png', dpi=150)
print("Saved benchmark_plots.png")