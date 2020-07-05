import pandas as pd
import numpy as np

if __name__=="__main__":
    data=pd.read_csv('./companylist.csv') # read the original csv file
    res=data[['Symbol','LastSale']] # extract the columns Symbol and LastSale
    res.dropna(inplace=True) # drop the NaN row
    res=np.array(res).tolist()
    for i in range(len(res)): # process the float numbers
        res[i][1]=res[i][1]*10000
        res[i][1]=int(res[i][1])
    
    # generate the new csv file
    name=['Symbol','LastSale']
    fin=pd.DataFrame(columns=name,data=res)
    fin.to_csv("preprocessedList.csv",encoding='utf-8')