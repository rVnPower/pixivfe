# this temporary script converts .csv to .yml
# can delete when documentation migration is done
import sys,csv
r = csv.reader(open(sys.argv[1]))
keys = next(r)
for row in r:
  print('-')
  for k, v in zip(keys, row):
    print(f'''  {k}: {v}''')
