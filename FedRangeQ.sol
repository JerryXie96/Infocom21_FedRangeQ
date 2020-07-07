pragma solidity = 0.5.4;
pragma experimental ABIEncoderV2;

contract FedRangeQ{
    // should be modified if the number of bits in one block changes
    struct IndexStru{               // the structure of one index item
        uint label;                 // the label of one index
        bytes[16] pubKey;           // the public key of each block [the number of blocks]
        uint8[3][3][16] tag;        // the ciphertext tag of each variable's ciphertext (sub-index) [max number of conflicts in one baskets][the number of baskets][the number of blocks]
        bytes[3][16] ciphertext;    // ciphertext of each block [the number of variables in one block][the number of blocks]
    }
    
    // should be modified if the number of bits in one block changes
    struct QueryStruct{
        bytes pubKey;           // the public key of one query
        uint8[16] tag;          // the ciphertext tag of each blocks [the number of blocks]
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
    
    // return the matching result
    function getResult() public view returns(uint[] memory){
        return result;
    }
    
    // clear the matching result
    function clearResult() public{
        delete result;
    }
    
}