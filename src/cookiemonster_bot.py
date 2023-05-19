import discord
from discord.ext import commands
from os import getenv

TOKEN = getenv("COOKIEMONSTER_TOKEN")
bot = commands.Bot(command_prefix="<", intents=discord.Intents.all())

@bot.event
async def on_ready():
    print(f"{bot.user} is ready!")

if __name__ == "__main__":
    bot.run(TOKEN)
