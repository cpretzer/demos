const express = require('express')
const hbs = require('express-hbs')
const path = require('path')
const app = express()
const {genExercises} = require('./dist/exercise')
const port = 3000

app.set('views', path.join(__dirname, 'views'))

app.set('view engine', 'hbs')

let exercises = genExercises()

// console.log(exercises)

app.get('/', (req, res) => {
  res.render('index', {exercises: exercises})
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})