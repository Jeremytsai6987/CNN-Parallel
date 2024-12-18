#!/bin/bash
#
#SBATCH --mail-user=jeremyyawei@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj1_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/jeremyyawei/Parallel_Programming/project1/project-1-Jeremytsai6987/proj1/
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=04:00:00

# Stop script on first error
set -e

echo "Loading modules"
module load golang/1.19

cd ./editor 
echo "Cleaning up old results"
rm /home/jeremyyawei/Parallel_Programming/project1/project-1-Jeremytsai6987/proj1/benchmark/results.csv
rm -f speedup-*.png

PROGRAM_VERSIONS=("sequential" "parfiles" "parslices")
DATASETS=("small" "mixture" "big")
THREAD_COUNTS=(2 4 6 8 12)
RUNS=5

for DATASET in "${DATASETS[@]}"; do
    PROGRAM="sequential"
    BEST_TIME_SEQ=999999
    for ((i=1; i<=RUNS; i++)); do
        echo "Sequential run $i for dataset $DATASET"
        START_TIME=$(date +%s%N)
        go run editor.go $DATASET  # Use correct relative path to editor.go
        END_TIME=$(date +%s%N)
        ELAPSED_TIME=$((END_TIME - START_TIME))
        ELAPSED_TIME_SEC=$(echo "scale=6; $ELAPSED_TIME/1000000000" | bc)
        if (( $(echo "$ELAPSED_TIME_SEC < $BEST_TIME_SEQ" | bc -l) )); then
            BEST_TIME_SEQ=$ELAPSED_TIME_SEC
        fi
    done
    echo "$PROGRAM,$DATASET,1,$BEST_TIME_SEQ" >> /home/jeremyyawei/Parallel_Programming/project1/project-1-Jeremytsai6987/proj1/benchmark/results.csv

    for PROGRAM in parfiles parslices; do
        echo "Running dataset: $DATASET ($PROGRAM version)"
        for THREADS in "${THREAD_COUNTS[@]}"; do
            export GOMAXPROCS=$THREADS
            BEST_TIME=999999
            for ((i=1; i<=RUNS; i++)); do
                echo "$PROGRAM run $i with $THREADS threads for dataset $DATASET"
                START_TIME=$(date +%s%N)
                go run editor.go $DATASET $PROGRAM $THREADS  # Use correct relative path to editor.go
                END_TIME=$(date +%s%N)
                ELAPSED_TIME=$((END_TIME - START_TIME))
                ELAPSED_TIME_SEC=$(echo "scale=6; $ELAPSED_TIME/1000000000" | bc)
                if (( $(echo "$ELAPSED_TIME_SEC < $BEST_TIME" | bc -l) )); then
                    BEST_TIME=$ELAPSED_TIME_SEC
                fi
            done
            echo "$PROGRAM,$DATASET,$THREADS,$BEST_TIME" >> /home/jeremyyawei/Parallel_Programming/project1/project-1-Jeremytsai6987/proj1/benchmark/results.csv
        done
    done
done

echo "Running Python plotting script"
python3 /home/jeremyyawei/Parallel_Programming/project1/project-1-Jeremytsai6987/proj1/benchmark/plot_speed_up.py

echo "Benchmarking script completed"
