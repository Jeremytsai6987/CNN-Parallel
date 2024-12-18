import matplotlib
matplotlib.use('Agg')
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

# Load data
df = pd.read_csv('/home/jeremyyawei/Parallel_Programming/project1/project-1-Jeremytsai6987/proj1/benchmark/results.csv', names=['Program', 'Dataset', 'Threads', 'BestTime'])

# Define colors for each dataset
colors = {'small': 'blue', 'mixture': 'green', 'big': 'orange'}

# Thread counts used in the benchmark
thread_counts = [2, 4, 6, 8, 12]

# Parallel programs to analyze
parallel_programs = ['parfiles', 'parslices']

# Fixed parallel fraction P
P = 0.95

# Function to calculate Amdahl's speedup
def amdahl_speedup(P, N):
    return 1 / ((1 - P) + (P / N))

# Iterate over each parallel program
for program in parallel_programs:
    plt.figure(figsize=(12, 10))  
    max_y_value = 0

    for dataset in ['small', 'mixture', 'big']:
        try:
            seq_time = df[(df['Dataset'] == dataset) & (df['Program'] == 'sequential')]['BestTime'].values[0]
        except IndexError:
            print(f"Sequential time not found for dataset '{dataset}'. Skipping.")
            continue
        
        subset = df[(df['Dataset'] == dataset) & (df['Program'] == program) & (df['Threads'].isin(thread_counts))].drop_duplicates(subset=['Threads'])
        
        if subset.empty:
            print(f"No data found for program '{program}' and dataset '{dataset}'. Skipping.")
            continue

        speedup_values = [seq_time / subset[subset['Threads'] == N]['BestTime'].values[0] for N in thread_counts]
        
        theoretical_speedup = [amdahl_speedup(P, N) for N in thread_counts]

        plt.plot(thread_counts, speedup_values, marker='o', color=colors[dataset], label=f'{dataset} measured', linewidth=2)
        
        plt.plot(thread_counts, theoretical_speedup, linestyle='--', color=colors[dataset], label='theoretical' if dataset == 'small' else None)

        #max_speedup_infinite = 1 / (1 - P)
        #plt.axhline(max_speedup_infinite, color='red', linestyle=':', label='Max theoretical (∞ threads)' if dataset == 'small' else None)

        max_y_value = max(max_y_value, max(speedup_values), max(theoretical_speedup))

    plt.title(f"Speedup Graph for '{program}' (Measured vs. Theoretical)")
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.xticks(thread_counts)
    plt.grid(True)
    
    plt.legend(title='Dataset', loc='upper left')
    
    plt.ylim(0, max_y_value * 1.2)  

    plt.tight_layout(pad=2.0)
    
    plt.savefig(f"proj1/benchmark/speedup-{program}.png")
    plt.close()
