# Mathcord: A Disord Math Expression Bot Written From Scratch

Good day! You can add a version of the bot via:

[<img height="200" width="200" src="https://avatars.githubusercontent.com/u/56949250?v=4" />](https://discord.com/api/oauth2/authorize?client_id=941278523132358686&scope=bot&permissions=380104611840)


Use:

```
/solve <expression>
```

This bot is written from scratch. It relies on zero third-party libraries (save for Godotenv).

The following algorithms have been implemented for this bot:

- Shunting Yard (for actual parsing and solving)
- SHA-512 (which ED25519 relies on)
- ED25519 (for verification)


To run this bot on your own server, just do:

<pre prefix="$">sudo docker-compose up -d</pre>

