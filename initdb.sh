#!/bin/bash
mongo --eval 'db.plugins.createIndex( { name: "text", description: "text" } )' plugins
mongo --eval 'db.sessions.createIndex( { "createdAt": 1 }, { expireAfterSeconds: 86400 } )' users