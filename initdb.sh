#!/bin/sh
mongo --eval 'db.plugins.createIndex( { name: 1 }, { unique: true } )' plugins
mongo --eval 'db.plugins.createIndex( { name: "text", description: "text" } )' plugins
mongo --eval 'db.sessions.createIndex( { createdAt: 1 }, { expireAfterSeconds: 86400 } )' users
mongo --eval 'db.users.createIndex( { username: 1 }, { unique: true } )' users
mongo --eval 'db.users.createIndex( { email: 1 }, { unique: true } )' users