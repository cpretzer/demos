const express = require('express')
// const http = require('http')
// const fetch = require('node-fetch')
const axios = require('axios')
const app = express()

process.argv.forEach((e, i) => {
    console.log("element is " + e + " at index " + i)
})

let serverType = process.argv[2]
let port = process.argv[3]
let remoteUrl

if (port == undefined) {
    port = 3000
}

if (serverType == undefined) {
    serverType = 'terminus'
} else if (serverType === 'point-to-point-channel') {
    remoteUrl = process.argv[4]

    if (remoteUrl === undefined || remoteUrl === '') {
        console.log('no remoteUrl set for point-to-point-channel type')
        process.exit(1)
    }

    console.log('remoteUrl is ' + remoteUrl)
}

console.log('serverType is ' + serverType)

app.get('/', async (req, res) => {
    if (serverType === 'terminus') {        
        console.log('sending terminus response Hello World')
        res.send('Hello World!')
    } else {
        console.log("making passthrough request")
        var promise = await axios.get(remoteUrl)
            .then(data => {
                console.log(data.data)
                res.send(data.data);
            })
            .catch(err => {
                console.log(err)
                res.send('unable to handle request')
            })
            .then( console.log('all done'))
    }
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})