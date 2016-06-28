CREATE DATABASE lcs DEFAULT CHARACTER SET = utf8;

USE lcs;

CREATE TABLE teams (
	teamid INTEGER AUTO_INCREMENT PRIMARY KEY,
	lcsid INTEGER NOT NULL,
	riotid INTEGER NOT NULL,
	name VARCHAR(40) NOT NULL,
	tag VARCHAR(10) NOT NULL
) ENGINE = InnoDB;

CREATE TABLE players (
	playerid INTEGER AUTO_INCREMENT PRIMARY KEY,
	lcsid INTEGER NOT NULL,
	riotid INTEGER NOT NULL,
	name VARCHAR(40) NOT NULL,
	teamid INTEGER NOT NULL,
	position INTEGER NOT NULL,
	addlpos TEXT,
	FOREIGN KEY (`teamid`) REFERENCES teams(`teamid`)
) ENGINE = InnoDB;

CREATE TABLE games (
	gameid INTEGER AUTO_INCREMENT PRIMARY KEY,
	lcsid INTEGER NOT NULL,
	riotgameid VARCHAR(40) NOT NULL,
	riotmatchid VARCHAR(40) NOT NULL,
	redteamid INTEGER NOT NULL,
	blueteamid INTEGER NOT NULL,
	gamestart DATETIME NOT NULL,
	gameend DATETIME DEFAULT NULL,
	FOREIGN KEY (`redteamid`) REFERENCES teams(`teamid`),
	FOREIGN KEY (`blueteamid`) REFERENCES teams(`teamid`)
) ENGINE = InnoDB;

CREATE TABLE stats (
	statid INTEGER AUTO_INCREMENT PRIMARY KEY,
	riotname VARCHAR(40) NOT NULL
) ENGINE = InnoDB;

CREATE TABLE updates (
	updateid BIGINT AUTO_INCREMENT PRIMARY KEY,
	tag VARCHAR(40) NOT NULL,
	time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

CREATE TABLE versions (
	tablename VARCHAR(20) PRIMARY KEY,
	version INTEGER
) ENGINE = InnoDB;
