export function isPromise(value) {
  if (value !== null && typeof value === 'object') {
    return value.promise && typeof value.promise.then === 'function';
  }
}

export function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}
export const checkIpFormat = (ip) => {
    let re = /(\d+)\.(\d+)\.(\d+)\.(\d+)/g
    return re.test(ip)
}
export const checkHostName = (str) => {
    let re = /^cghost[0-9]*[1-9][0-9]*$/g
    return re.test(str)
}

export const Atype_Zone = 1
export const Atype_ZoneDB = 2
export const Atype_ZoneLogDB = 3
export const Atype_ZONEBAKDB = 4
export const Atype_DATALOGDB = 5
export const Atype_Frame = 6

export const checkAppliactionType = (str) => {
    let re = /^zone[0-9]*[1-9][0-9]*$/g
    if (re.test(str)) {
        return Atype_Zone
    }
    let re1 = /^zonedb[0-9]*[1-9][0-9]*$/g
    if (re1.test(str)) {
        return Atype_ZoneDB
    }
    let re2 = /^zonelogdb[0-9]*[1-9][0-9]*$/g
    if (re2.test(str)) {
        return Atype_ZoneLogDB
    }
    let re3 = /^zonedbBak[0-9]*[1-9][0-9]*$/g
    if (re3.test(str)) {
        return Atype_ZONEBAKDB 
    }
    let re4 = /^datalogdb[0-9]*[1-9][0-9]*$/g
    if (re3.test(str)) {
        return Atype_DATALOGDB 
    }

    if (str == ""){
        return 0
    }

    return Atype_Frame
}

export const trim = (str) => str.replace(/(^\s+)|(\s+$)/g, "")

