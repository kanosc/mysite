const copyList = function (a) {
    return JSON.parse(JSON.stringify(a))
}
const convertSelectToColor = (key) => {
    const m = new Map([
        ["Dark", "#996633"],
        ["Red", "#FF0000"],
        ["Yellow", "#FFFF33"],
        ["Green", "#48D1CC"],
        ["Orange", "#fda500"],
        ["Blue", "#99CCFF"],
        ["Gray", "#999999"],
        ["Pink", "#C71585"]
    ])
    return m.get(key)
}
export { copyList, convertSelectToColor }