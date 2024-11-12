import argparse
import csv
import os
import random

# Set up argument parser
parser = argparse.ArgumentParser(description='Generate random CSV data')
parser.add_argument('num', help='Number of rows to generate')
parser.add_argument('output_file', help='Output CSV filename')
args = parser.parse_args()

# Number of rows to generate
num_rows = int(args.num)

# Generate random data
data = []
for _ in range(num_rows):
    filepath = f"/path/to/file_{random.randint(1, 1000)}.txt"
    float_val = random.uniform(0, 40.0)
    uint_val = random.randint(0, 2000)
    data.append([filepath, float_val, uint_val])

# Write to CSV file
with open(args.output_file, 'w', newline='') as f:
    writer = csv.writer(f)
    # Write data rows
    writer.writerows(data)
