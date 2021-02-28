import requests

""" Script Inputs """
application_id = ""
guild_id = ""  # aka server ID
bot_token = ""
command_id = ""
create_command_body = {
    # See https://discord.com/developers/docs/interactions/slash-commands#create-global-application-command
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

headers = {
    "Authorization": f"Bot {bot_token}"
}


def get_commands(guild_id=None):
    print("Getting commands to see if it has been updated")
    if guild_id != None:
        url = f"https://discord.com/api/v8/applications/{application_id}/guilds/{guild_id}/commands"
    else:
        url = f"https://discord.com/api/v8/applications/{application_id}/commands"
    r = requests.get(url, headers=headers)
    print(r.status_code)
    print(r.json())


def create_command(guild_id=None):
    if guild_id != None:
        url = f"https://discord.com/api/v8/applications/{application_id}/guilds/{guild_id}/commands"
    else:
        url = f"https://discord.com/api/v8/applications/{application_id}/commands"
    r = requests.post(url, headers=headers, json=create_command_body)
    print(r.status_code)
    print(r.json())


def delete_command(guild_id=None):
    if guild_id != None:
        url = f"https://discord.com/api/v8/applications/{application_id}/guilds/{guild_id}/commands/{command_id}"
    else:
        url = f"https://discord.com/api/v8/applications/{application_id}/commands/{command_id}"
    r = requests.delete(url, headers=headers)
    print(r.status_code)


if __name__ == "__main__":
    # create_command()
    get_commands()
    # delete_command()
