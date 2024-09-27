DROP DATABASE IF EXISTS GuessingGameDb;

CREATE DATABASE GuessingGameDb;

USE GuessingGameDb;

CREATE TABLE Users (
  Id int NOT NULL Auto_Increment,
  Email varchar(200) UNIQUE NOT NULL,
  PassCode varchar(60) NOT NULL,
  SettingsGuess int,
  SettingsBox int,
  PRIMARY KEY(Id)
);

CREATE TABLE Games (
    Id int NOT NULL AUTO_INCREMENT,
    Guesses int NOT NULL,
    Result boolean,
    GamesWon int,
    GamesLost int,
    MaxGuesses int, 
    UsersId int,
    PRIMARY KEY (Id),
    FOREIGN KEY(UsersId) REFERENCES Users(Id)
);