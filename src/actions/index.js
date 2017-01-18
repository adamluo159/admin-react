import fetch from 'isomorphic-fetch'

const actionCreator = (regActionType,type)=>{
   regActionType[type] = type
   return (playload)=>({
        type,
        playload
   }
)}

const fetchMachines = ()=>{
    return dispath=>{
        dispath(machine.reqMachines())
        let data = [{
               key: "host0",
               hostname: 'host0',
               IP: "192.168.1.1",
               outIP:"192.168.1.1",
               type :"login",
               edit: false,
        }]
        dispath(machine.recvMachines({
            data:data
        }))
       // return fetch('')
       // .then(response=>response.json())
       // .then(json=>dispath(RecvLayoutData(loading, json)))
    }
}

export const actions={}
const  machine={
    "selectMainlayout": actionCreator(actions, 'SELECT_MAINLAYOUT_KEY'),
    "reqMachines":      actionCreator(actions, 'REQ_MACHINES'),
    "recvMachines":     actionCreator(actions, 'RECV_MACHINES'),
    "addMachine":       actionCreator(actions, 'ADD_MACHINE'),
    "editMachine":      actionCreator(actions, 'EDIT_MACHINE'),
    "saveMachine":      actionCreator(actions, 'SAVE_MACHINE'),
    "delMachine" :      actionCreator(actions, 'DEL_MACHINE'),
    "pageMachine" :     actionCreator(actions, 'PAGE_MACHINE'),
    "fetchMachines":    fetchMachines,
}
export default machine
