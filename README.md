# BeanBot

Can you count in a Discord text channel? Probably...

BeanBot lets you designate a text channel as the "Learn to Count" channel. It will keep track of you counting in that channel like so

```
alice: 1
bob: 2
alice: 3
josh: 4
bob: 5
eve: 6
```

Every time you add a number to the count, you will receive that amount of beans! But, if you make a mistake, or count twice in a row...

```
alice: 7
eve: 7
BeanCounter: Eve spilled the beans and lost 28 beans!
```

You lose all of the beans that have been given out since you started! This can lead to some very large penalties. Making a mistake at 200 beans in will lose you 20,100 beans! It's suprisingly easy to mess up the count.

You can add BeanBot to your server with [this link](https://discord.com/api/oauth2/authorize?client_id=759168728554012672&permissions=75776&scope=bot). To set it up, type `!bean configure` in your prospective counting channel (as someone with the Manage Channels permission).

# Features

- Tracks counting in a channel
- `!bean query` to check your progress
- `!bean leaderboard` to see who's the best at counting
- `!bean give` to create some sick, twisted economy based on this
- `!bean bet` to further your gambling addiction

There are some various other features which can be seen with just `!bean` or `!bean help`.

# Contributing

Feel free to open a PR. Please `go fmt` your code before doing so and ensure your code at least compiles. This repo is tagged for Hacktoberfest, because I was frustrated at how many pseudo-projects popped up to make faux PRs for that event, and wanted to be part of the solution instead of just complaining.

I don't really have specific features in mind for the bot - if there are minor things I notice to be fixed they can be found on the issues page.
