import matplotlib.pyplot as plt
import matplotlib.patches as mpatches

# Events: (start_minutes, duration_minutes, label, color)
# Reference start = 0 at first attempt
events = [
    (0, 4,   '2.0 (failed)', 'red'),     # 4 minutes
    (4, 4,   '2.1 (failed)', 'orange'),  # 4 minutes
    (8, 29,  '2.2 (failed)', 'yellow'),  # 29 minutes
    (37, 2,  '2.4 (success)', 'green'),  # 2 minutes
]

fig, ax = plt.subplots(figsize=(10,3))
for i, (start, dur, label, color) in enumerate(events):
    ax.barh(i, dur, left=start, color=color, edgecolor='black')
    ax.text(start + dur/2, i, label, ha='center', va='center', fontsize=10, color='black')

ax.set_yticks(range(len(events)))
ax.set_yticklabels([f'Attempt {i+1}' for i in range(len(events))])
ax.set_xlabel('Time (minutes from first attempt)')
ax.set_title('Chaincode Deployment Attempts (M4)')
plt.tight_layout()
plt.savefig('deployment_timeline.png', dpi=150)
# No plt.show() to avoid backend warning