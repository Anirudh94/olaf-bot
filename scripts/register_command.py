import requests

""" Script Inputs """
application_id = ""
guild_id = None  # aka server ID
bot_token = ""
command_id = ""
""" End script inputs """

create_command_body = {
    # See https://discord.com/developers/docs/interactions/slash-commands#create-global-application-command
    "name": "image",
    "description": "Send a random photo",
    "options": [
        {
            "name": "keywords",
            "description": "Keywords of the photo",
            "type": 3,
            "required": True,
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
    # create_command(guild_id=guild_id)
    get_commands(guild_id=guild_id)
    # delete_command(guild_id=guild_id)
