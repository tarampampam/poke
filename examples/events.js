#!/usr/bin/env poke run

events.push({level: 'info', message: 'This is an info message'})
events.push({level: 'warning', message: 'This is a warning message', error: new Error('optional error')})
// events.push({level: 'error', message: 'This is a error message', error: new Error('optional error')})
