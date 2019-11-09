import json
import requests
import random

file_out = open('test.sh', 'w')

url = "http://127.0.0.1:55555/set_board_data"

headers = {
    'Content-Type': "application/json",
    'cache-control': "no-cache",
    'Postman-Token': "7bde98f3-9714-4c23-aa45-d39708e7f13d"
    }

for ind in range(1,100000):
    last = random.randint(0,10)
    mac = '00:00:00:00:'+(str(last) if len(str(last))==2 else '0'+str(last))
    temperature = random.randint(-50,50)
    
    Types = ['temperatures', 'humidity', 'pressure', 'soil']
    Type = Types[random.randint(0,3)]


    Obj = {"mac" : mac, "valueType" : Type, "value" : temperature, "unit" : "C"}
    
    sting = r"curl -X POST http://127.0.0.1:55555/set_board_data -H 'Content-Type: application/json' -H 'Postman-Token: db9079db-a701-4154-aa1e-8232fa6a0a89' -H 'cache-control: no-cache' -d '"+json.dumps(Obj)+"' &"
    file_out.write(sting+'\n')
    # file_out.write('sleep '+str(random.randint(0,60))+'s\n')

file_out.close()

    # Obj = {"mac" : mac, "valueType" : "temperature", "value" : temperature, "unit" : "C"}
    # print(Obj)
    # response = requests.request('POST', url, data = Obj, headers=headers )
    # print(response.url)
    # if response.status_code == 201:
    #     print(response.content)
    # else:
    #     print(response)
    #     exit(1)
    