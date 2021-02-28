import requests

application_id=""
guild_id=""
bot_token=""
url = f"https://discord.com/api/v8/applications/{application_id}/guilds/{guild_id}/commands"
headers = {
    "Authorization": f"Bot {bot_token}"
}

# See https://discord.com/developers/docs/interactions/slash-commands#create-global-application-command
json = {
    "name": "blep",
    "description": "Send a random adorable animal photo",
    "options": [
        {
            "name": "animal",
            "description": "The type of animal",
            "type": 3,
            "required": True,
            "choices": [
                {
                    "name": "Dog",
                    "value": "animal_dog"
                },
                {
                    "name": "Cat",
                    "value": "animal_cat"
                },
                {
                    "name": "Penguin",
                    "value": "animal_penguin"
                }
            ]
        },
        {
            "name": "only_smol",
            "description": "Whether to show only baby animals",
            "type": 5,
            "required": False
        }
    ]
}

r = requests.post(url, headers=headers, json=json)
print(r.status_code)
print(r.json())

print("Getting commands to see if it has been updated")
r = requests.get(url, headers=headers)
print(r.status_code)
print(r.json())