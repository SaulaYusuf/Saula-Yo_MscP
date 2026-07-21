import matplotlib.pyplot as plt

# Real counts from `grep -c "^func" main.go` for each version
versions = ['1.0', '2.0', '2.4']
functions = [2, 4, 4]  # 1.0 had RecordTelemetry, ReadTwin; 2.0 added RecordHandover, ReadHandover; 2.4 same

plt.figure(figsize=(6,4))
plt.plot(versions, functions, marker='o', linestyle='-', color='purple')
plt.xlabel('Chaincode Version')
plt.ylabel('Number of Public Functions')
plt.title('Function Count per Version')
for i, val in enumerate(functions):
    plt.text(versions[i], val + 0.1, f'{val}', ha='center')
plt.tight_layout()
plt.savefig('function_count.png', dpi=150)
plt.show()