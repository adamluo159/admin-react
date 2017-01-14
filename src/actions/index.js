import fetch from 'isomorphic-fetch'
export const REQ_LAYOUT_POSTS  = 'REQ_LAYOUT_POSTS'
export const RECV_LAYOUT_POSTS = 'RECV_LAYOUT_POSTS'
export const SELECT_MAINLAYOUT_KEY='SELECT_MAINLAYOUT_KEY'

export const selectMainlayout=(key)=>{
    console.log(key)
    return {
        type: SELECT_MAINLAYOUT_KEY,
        selectKey: key
    }
} 

export const ReqLayoutData = (loading)=>{
    return {
        type: REQ_LAYOUT_POSTS,
        loading,
    }
}

export const GetLayoutData = (loading)=>{
    return dispath=>{
        dispath(ReqLayoutData())
        return fetch('')
        .then(response=>response.json())
        .then(json=>dispath(RecvLayoutData(loading, json)))
    }
}

export const RecvLayoutData = (loading, jsonData)=>{
    return {
        type: RECV_LAYOUT_POSTS,
        loading,
        jsonData,
    }
}