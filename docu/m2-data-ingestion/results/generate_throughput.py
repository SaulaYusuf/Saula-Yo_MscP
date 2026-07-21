import matplotlib.pyplot as plt

# Real data
records = 8000
time_sec = 123.8
tps = records / time_sec

# Baseline comparison: theoretical max (e.g., 100 TPS) for visual context
theoretical_max = 100

categories = ['Actual TPS', 'Theoretical Max']
values = [tps, theoretical_max]

plt.figure(figsize=(6,4))
bars = plt.bar(categories, values, color=['#1f77b4', '#ff7f0e'])
plt.ylabel('Transactions per Second (TPS)')
plt.title('Sensor Ingestion Throughput (M2)')
plt.ylim(0, 120)
for bar, val in zip(bars, values):
    plt.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 2, f'{val:.2f}', ha='center')
plt.tight_layout()
plt.savefig('throughput.png', dpi=150)