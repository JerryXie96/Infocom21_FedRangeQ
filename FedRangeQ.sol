pragma solidity = 0.5.4;
pragma experimental ABIEncoderV2;

contract FedRangeQ{
    // should be modified if the number of bits in one block changes
    struct IndexStru{               // the structure of one index item
        uint label;                 // the label of one index
        // next line should be changed if bit length changes
        bytes[16] pubKey;           // the public key of each block [the number of blocks]
        // next line should be changed if bit length changes
        uint8[3][3][16] tag;        // the id of ciphertext with the same tag, the unused items will be filled with 100 (i.e., a value out of the range). e.g. tag[0][1][] includes all the ids in block 0 which their tags are 1 (sub-index) [max number of conflicts in one baskets][the number of baskets][the number of blocks]
        // next line should be changed if bit length changes
        bytes[3][16] ciphertext;    // ciphertext of each block [the number of variables in one block][the number of blocks]
    }
    
    // should be modified if the number of bits in one block changes
    struct QueryStru{
        bytes pubKey;           // the public key of one query
        // next line should be changed if bit length changes
        uint8[16] tag;          // the ciphertext tag of each blocks [the number of blocks]
        // next line should be changed if bit length changes
        bytes[16] ciphertext;   // the ciphertext of each block [the number of blocks]
    }
    
    IndexStru[] index;  // the index in blockchain
    uint[] result;     // the matching result
    
    // store the uploaded index in blockchain, indexToBeAdded is a list which includes some index items
    function store(IndexStru[] memory indexToBeAdded) public {
        for(uint i=0;i<indexToBeAdded.length;i++){ // store the index one by one
            index.push(indexToBeAdded[i]);
        }
    }
    
    // should be modified if the number of bits in one block changes
    // search the matched items from the index
    function search(QueryStru memory query) public {
        for(uint i=0;i<index.length;i++){   // scan each item in the index
            // next line should be changed if bit length changes
            bool isMatched=false;           // the flag which shows whether the current item in the index is matched
            for(uint j=0;j<16;j++){         // scan each block in query
                // next line should be changed if bit length changes
                for(uint k=0;k<3;k++){      // scan the tag list of the current block which their tags are the same as query's
                    if(index[i].tag[j][query.tag[j]][k]==100){  // if all the blocks which their tags are the same as query's are checked, break the loop
                        break;       
                    }
                    uint8 id=index[i].tag[j][query.tag[j]][k];   // get the id in ciphertext array
                    
                    // use the bn256 bilinear map
                    bytes memory input=abi.encodePacked(index[i].ciphertext[j][id],query.pubKey,index[i].pubKey[j],query.ciphertext[j]);    // pack the bytes data
                    uint[1] memory output;  // to store the bilinear result
                    uint length=input.length;
                    assembly{
                        if iszero(call(not(0),0x08,0,add(input,0x20),length,output,0x20)){  // perform the bilinear map
                            revert(0,0)
                        }
                    }
                    if(output[0]!=0){       // the ciphertext is matched
                        isMatched=true;     // the value is matched
                        break;
                    }
                }
                if(isMatched==true){        // the value is matched
                    break;
                }
            }
            if(isMatched==true){
                result.push(index[i].label);    // add the index into the result list
            }
        }
    }
    
    // return the matching result
    function getResult() public view returns(uint[] memory){
        return result;
    }
    
    // clear the matching result
    function clearResult() public{
        delete result;
    }
    
}