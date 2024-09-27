const express = require('express');
const app = express();
const port = 3000;
const bcrypt = require('bcryptjs');
const http = require('http');
const fs = require('fs');

const conInfo =
{
  host: 'localhost',
  user: 'root',
  password: '',
  database: 'GuessingGameDb'
};

const mysql = require('mysql');
const session = require('express-session');
app.use(session({
  secret: 'happy jungle',
  resave: false,
  saveUninitialized: false,
  cookie: { maxAge: 600000 }
}));

app.all('/', serveIndex);
app.all('/whoIsLoggedIn', whoIsLoggedIn);
app.all('/register', register);
app.all('/login', login);
app.all('/logout', logout);
app.all('/game', game);
app.all('/stats', stats);
app.all('/history', history);
app.all('/updateUser', updateUser);
app.listen(port, 'localhost', startHandler());

function startHandler() {
  console.log('Server listening on port ' + port);
}

function register(req, res) {
  if (req.query.email == undefined || !validateEmail(req.query.email)) {
    writeResult(req, res, { 'error': 'Please enter a valid email.' });
    return;
  }

  if (req.query.password == undefined || !validatePassword(req.query.password)) {
    writeResult(req, res, { 'error': 'Password should be between 8-20 characters.' });
    return;
  }

  if (req.query.boxCount == undefined || req.query.maxGuesses == undefined) {
    writeResult(req, res, { 'error': 'Must set number of boxes and max attempts allowed.' });
    return;
  }

  if (parseInt(req.query.boxCount) < parseInt(req.query.maxGuesses)) {
    writeResult(req, res, { 'error': 'Max attempts must be less than the number of boxes.' });
    return;
  }

  if (isNaN(req.query.boxCount) || isNaN(req.query.maxGuesses)) {
    writeResult(req, res, { 'error': 'Number of boxes and max attempts must be valid numbers.' });
    return;
  }

  let con = mysql.createConnection(conInfo);
  con.connect(function (err) {
    if (err)
      writeResult(req, res, { 'error': err });
    else {
      let hash = bcrypt.hashSync(req.query.password, 12);
      con.query('INSERT INTO Users (Email, PassCode, SettingsBox, SettingsGuess) VALUES (?, ?, ?, ?)', [req.query.email, hash, req.query.boxCount, req.query.maxGuesses], function (err, result, fields) {
        if (err) {
          if (err.code == 'ER_DUP_ENTRY')
            err = 'User account already exists.';
          writeResult(req, res, { 'error': err });
        }
        else {
          con.query('SELECT * FROM Users WHERE Email = ?', [req.query.email], function (err, result, fields) {
            if (err)
              writeResult(req, res, { 'error': err });
            else {
              req.session.user = { 'id': result[0].Id, 'email': result[0].Email, 'userBoxCount' : result[0].SettingsBox, 'userMaxGuesses' : result[0].SettingsGuess };
              req.session.userId = result[0].Id;
              writeResult(req, res, { 'user': req.session.user });
            }
          });
        }
      });
    }
  });
}

function login(req, res) {
  if (req.query.email == undefined) {
    writeResult(req, res, { 'error': "Email is required" });
    return;
  }

  if (req.query.password == undefined) {
    writeResult(req, res, { 'error': "Password is required" });
    return;
  }

  let con = mysql.createConnection(conInfo);
  con.connect(function (err) {
    if (err)
      writeResult(req, res, { 'error': err });
    else {
      con.query("SELECT * FROM Users WHERE Email = ?", [req.query.email], function (err, result, fields) {
        if (err)
          writeResult(req, res, { 'error': err });
        else {
          if (result.length == 1 && bcrypt.compareSync(req.query.password, result[0].PassCode)) {
            req.session.user = { 'id': result[0].Id, 'email': result[0].Email, 'userBoxCount' : result[0].SettingsBox, 'userMaxGuesses' : result[0].SettingsGuess };
            req.session.userId = result[0].Id;
            writeResult(req, res, { 'user': req.session.user });
          }
          else {
            writeResult(req, res, { 'error': "Invalid email/password" });
          }
        }
      });
    }
  });
}

function logout(req, res) {
  req.session.user = undefined;
  writeResult(req, res, { 'nobody': 'Nobody is logged in.' });
}

function whoIsLoggedIn(req, res) {
  if (req.session.user == undefined)
    writeResult(req, res, { 'nobody': 'Nobody is logged in.' });
  else
    writeResult(req, res, { 'user': req.session.user });
}

function game(req, res) {
  // if we have not picked a secret number, restart the game...
  if (req.session.user == undefined)
    writeResult(req, res, { 'nobody': 'Nobody is logged in.' });
  else {
    if (req.session.answer == undefined) {
      req.session.guesses = 0;
      req.session.maxGuesses = req.query.maxGuesses;
      req.session.answer = Math.floor(Math.random() * req.query.boxCount) + 1;
      console.log(req.session.answer);
    }

    // if a guess was not made, restart the game...
    if (req.query.guess == undefined) {
      req.session.answer = undefined;
      req.session.guesses = 0;
      req.session.maxGuesses = req.query.maxGuesses;
      req.session.answer = Math.floor(Math.random() * req.query.boxCount) + 1;
      console.log(req.session.answer);
      writeResult(req, res, { 'gameStatus': 'Start', 'user': req.session.user });
    }
    else {
      // a guess was made, check to see if it is correct...
      if (req.query.guess == req.session.answer) {
        req.session.guesses = req.session.guesses + 1;
        let con = mysql.createConnection(conInfo);
        con.connect(function (err) {
          if (err)
            console.log(err);
          else {
            con.query('INSERT INTO Games(Guesses, Result, GamesWon, MaxGuesses, UsersId) VALUES(?, ?, ?, ?, ?)', [req.session.guesses, true, 1, req.session.maxGuesses, req.session.userId], function (err, result, fields) {
              if (err)
                console.log(err);
              else
                console.log('Inserted Game Record.');
            });
          }
        });
        req.session.answer = undefined;
        writeResult(req, res, { 'gameStatus': 'Win', 'guesses': req.session.guesses, 'user': req.session.user });
      }
      else if (req.session.guesses + 1 == req.session.maxGuesses) {
        req.session.guesses = req.session.guesses + 1;
        let con = mysql.createConnection(conInfo);
        con.connect(function (err) {
          if (err)
            console.log(err);
          else {
            con.query('INSERT INTO Games(Guesses, Result, GamesLost, MaxGuesses, UsersId) VALUES(?, ?, ?, ?, ?)', [req.session.guesses, false, 1, req.session.maxGuesses, req.session.userId], function (err, result, fields) {
              if (err)
                console.log(err);
              else
                console.log('Too many guesses.');
            });
          }
        });
        req.session.answer = undefined;
        writeResult(req, res, { 'gameStatus': 'gameOver', 'guesses': req.session.guesses, 'user': req.session.user });
      }
      // a guess was made, check to see if too high...
      else if (req.query.guess > req.session.answer) {
        req.session.guesses = req.session.guesses + 1;
        console.log(req.session.guesses);
        writeResult(req, res, { 'gameStatus': 'Lose', 'guesses': req.session.guesses, 'user': req.session.user });
      }
      // a guess was made, it must be too low...
      else {
        req.session.guesses = req.session.guesses + 1;
        console.log(req.session.guesses);
        writeResult(req, res, { 'gameStatus': 'Lose', 'guesses': req.session.guesses, 'user': req.session.user });
      }
    }
  }
}

function stats(req, res) {
  if (req.session.user == undefined)
    writeResult(req, res, { 'nobody': 'Nobody is logged in.' });
  else {
    let con = mysql.createConnection(conInfo);
    con.connect(function (err) {
      if (err)
        writeResult(req, res, { 'error': err });
      else {
        con.query('SELECT COUNT(Id) AS gamesPlayed, COUNT(GamesWon) AS gamesWon, COUNT(GamesLost) AS gamesLost FROM Games WHERE UsersId = ?', [req.session.userId], function (err, result, fields) {
          if (err)
            writeResult(req, res, { 'error': err });
          else
            writeResult(req, res, { 'statPlayed': result[0].gamesPlayed, 'statWon' : result[0].gamesWon, 'statLost' : result[0].gamesLost,  'user' : req.session.user });
        });
      }
    });
  }
}

function history(req, res) {
  if (req.session.user == undefined)
    writeResult(req, res, { 'nobdoy': 'Nobody is logged in.' });
  else {
    let con = mysql.createConnection(conInfo);
    con.connect(function (err) {
      if (err)
        writeResult(req, res, { 'error': err });
      else {
        con.query('SELECT Result, Guesses, MaxGuesses FROM Games WHERE UsersId = ? ORDER BY Id DESC', [req.session.userId], function (err, result, fields) {
          if (err)
            writeResult(req, res, { 'error': err });
          else
            writeResult(req, res, { 'history': result, 'user' : req.session.user});
        });
      }
    });
  }
}

function updateUser(req, res){
  if (req.session.user == undefined)
    writeResult(req, res, { 'nobody': 'Nobody is logged in.' });
  else {
    let con = mysql.createConnection(conInfo);
    con.connect(function (err) {
      if (err)
        writeResult(req, res, { 'error': err });
      else {
        con.query('UPDATE Users SET SettingsBox = ?, SettingsGuess = ? WHERE Id = ?', [req.query.boxCount, req.query.maxGuesses, req.session.userId], function (err, result, fields) {
          if (err)
            writeResult(req, res, { 'error': err });
          else
            writeResult(req, res, { 'user' : req.session.user });
        });
      }
    });
  }
}

function serveIndex(req, res) {
  res.writeHead(200, { 'Content-Type': 'text/html' });
  let index = fs.readFileSync('index.html');
  res.end(index);
}

function writeResult(req, res, obj) {
  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.write(JSON.stringify(obj));
  res.end('');
}

function validateEmail(email) {
  if (email == undefined) {
    return false;
  }
  else {
    let re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
  }
}

function validatePassword(pass) {
  if (pass == undefined) {
    return false;
  }
  else {
    let re = /^(?=.*[A-Za-z])(?=.*[0-9])(?=.*\d)[A-Za-z0-9\d]{8,20}$/;
    return re.test(pass);
  }
}