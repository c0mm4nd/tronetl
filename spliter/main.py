#%%
import argparse
import os
import csv

parser = argparse.ArgumentParser(description='Spliter')
parser.add_argument('input', type=str, help="csv file to be split")
parser.add_argument('output', type=str, help="output csv filename's prefix")
parser.add_argument('step', type=int, help="split step")

args = parser.parse_args()

input = open(args.input)
reader = csv.DictReader(input)

output = None
writer = None
output_name = ""
output_size = 0
last_save = 0

for line in reader:
    # line = origin.readline()
    # if not line: break
    # i = line.index(",")
    # str_block_number = line[0:i]
    block_number = int(line["number"]) # int(str_block_number)
    # fix gas_limit & gas_used
    line["gas_limit"] = int(line["gas_limit"], 16)
    line["gas_used"] = int(line["gas_used"], 16)

    if output is None:
        output_name = "/mnt/nvme2t0/output/output/blocks-" + str(block_number)
        output = open(output_name , "w")
        writer = csv.DictWriter(output, reader.fieldnames)

    if block_number != last_save and block_number % args.step == 1:
        last_save = block_number
        if output is not None:
            output.close()
            os.rename(output_name, output_name + "-" + str(block_number -1) + ".csv")
        output_name = args.output + str(block_number)
        output = open(output_name , "w")
        writer = csv.DictWriter(output, reader.fieldnames)
        
    writer.writerow(line)
    # output_size += len(line)

if output is not None:
    output.close()
    os.rename(output_name, output_name + "-" + str(block_number) + ".csv")
