#!/usr/bin/env poke run

events.push({level: 'info', message: 'This is an info message'})
events.push({level: 'warning', message: 'This is a warning message', error: new Error('optional error')})
