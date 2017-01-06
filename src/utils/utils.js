export const checkIpFormat = (ip)=>{
    let re=/(\d+)\.(\d+)\.(\d+)\.(\d+)/g 
    return re.test(ip)
}

export const trim = (str)=>(
    console.log(str, typeof(str))
)