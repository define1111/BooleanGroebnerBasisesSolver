import subprocess
import itertools
import os
import tempfile

def generate(n: int) -> list:
  equations = []
  for i in range(n):
    for j in range(i+1, n):
      monoms = []
      for k in range(n):
        if k == i or k == j:
          continue
        monoms.append(f"x{k+1}")
      monoms.append(f"x{i+1}*x{j+1}")
      equations.append(' + '.join(monoms))
      monoms.append("1")
      equations.append(' + '.join(monoms))
  return equations
      
def run(n):
  print(f"Check n = {n}")
  all = generate(n)
  print('\n'.join(all)) 
  for count in range(2, 4):
    print(f"Choose {count} equations")
    for perm in itertools.permutations(all, count):
      temp_in = tempfile.NamedTemporaryFile(mode='w')
      #print(temp_in.name)
      temp_in.write(f"{n}\n")
      temp_in.write('\n'.join(perm))
      temp_in.flush()
      cmd = f'go run main.go {temp_in.name}'
      subprocess.run(cmd, stdout=output, stderr=subprocess.DEVNULL, shell=True)
      output.flush()
      separator = '#' * 30
      output.write(separator + '\n')
      output.flush()
      temp_in.close()
      #break
    print()

output_filename = 'output.txt'
output = open(output_filename, 'w')
run(3)
run(4)
output.close()
