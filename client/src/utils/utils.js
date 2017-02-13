export const checkIpFormat = (ip) => {
    let re = /(\d+)\.(\d+)\.(\d+)\.(\d+)/g
    return re.test(ip)
}

export const trim = (str) => (console.log(str, typeof(str)))

export const actionCreator = (regActionType, type) => {
    regActionType[type] = type
    return (playload) => ({type, playload})
}

export const create = (dispatchCreator, dispatch) => {
    let dispatchsObj = {}
    Object
        .keys(dispatchCreator)
        .forEach(item => {
            dispatchsObj[item] = (e) => dispatch(dispatchCreator[item](e))
        })
    return dispatchsObj
}
