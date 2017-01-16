import fetch from 'isomorphic-fetch'
export const SELECT_MAINLAYOUT_KEY='SELECT_MAINLAYOUT_KEY'

export const REQ_MACHINES  = 'REQ_MACHINES'
export const RECV_MACHINES = 'RECV_MACHINES'
export const ADD_MACHINE   = 'ADD_MACHINE'
export const EDIT_MACHINE  = 'EDIT_MACHINE'
export const SAVE_MACHINE  = 'SAVE_MACHINE'
export const DEL_MACHINE  =  'DEL_MACHINE'
export const PAGE_MACHINE = 'PAGE_MACHINE'

export const selectMainlayout=(key)=>{
    return {
        type: SELECT_MAINLAYOUT_KEY,
        selectKey: key
    }
} 

export const reqMachines = ()=>{
    return {
        type: REQ_MACHINES,
    }
}

export const recvMachines = (json)=>{
    return {
        type: RECV_MACHINES,
        data: json,
    }
}

export const fetchMachines = ()=>{
    return dispath=>{
        dispath(reqMachines())
        dispath(recvMachines(
            [{  
               key: "host0",
               hostname: 'host0',
               IP: "192.168.1.1",
               outIP:"192.168.1.1",
               type :"login",
               edit: false,
            }]
        ))
        //return fetch('')
        //.then(response=>response.json())
        //.then(json=>dispath(RecvLayoutData(loading, json)))
    }
}

export const addMachine = (newItem)=>{
    return {
        type: ADD_MACHINE,
        newItem,
    }
}

export const editMachine = (index)=>{
    return {
        type : EDIT_MACHINE,
        index,
    }
}

export const saveMachine = (save) =>{
    return {
        type : SAVE_MACHINE,
        save,
    }
}

export const delMachine = (index) =>{
    return {
        type: DEL_MACHINE,
        index,
    }
}

export const pageMachine = (pagination, filters, sorter)=>{
    return {
        type: PAGE_MACHINE,
        page:{
              current: pagination.current,
              pageSize:pagination.pageSize,
        }
    }
}

