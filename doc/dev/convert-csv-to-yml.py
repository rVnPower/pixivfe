import sys,csv
r = csv.reader(open(sys.argv[1]))
keys = next(r)
for row in r:
  print('-')
  for k, v in zip(keys, row):
    print(f'''  {k}: {v}''')
