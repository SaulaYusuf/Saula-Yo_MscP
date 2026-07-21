import matplotlib.pyplot as plt
import numpy as np

# Real data from your ingestion logs
datasets = ['Sensor (IoT)', 'Logistics', 'Metadata']
records = [8000, 1000, 500]
time_sec = [123.8, 8.7, 4.71]
tps = [records[i] / time_sec[i] for i in range(3)]

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(12, 5))

# Bar chart: Records
bars1 = ax1.bar(datasets, records, color=['#1f77b4', '#ff7f0e', '#2ca02c'])
ax1.set_ylabel('Number of Records')
ax1.set_title('Dataset Size')
for bar, val in zip(bars1, records):
    ax1.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 50, f'{val}', ha='center')

# Bar chart: TPS
bars2 = ax2.bar(datasets, tps, color=['#1f77b4', '#ff7f0e', '#2ca02c'])
ax2.set_ylabel('Transactions per Second (TPS)')
ax2.set_title('Ingestion Throughput')
for bar, val in zip(bars2, tps):
    ax2.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 2, f'{val:.1f}', ha='center')

plt.tight_layout()
plt.savefig('results/ingestion_comparison.png', dpi=150)
print("Saved ingestion_comparison.png")