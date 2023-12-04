import re
import click
import requests

class User:
    def __init__(self, username, password, roles):
        self.username = username
        self.password = password
        self.roles = roles


users = {
    'admin': User('admin', 'admin', ['admin']),
    'user': User('user', 'user_password', ['user']),
}

def authenticate(username, password):
    user = users.get(username)
    if user and user.password == password:
        return user
    return None

def check_role(user, required_role):
    return required_role in user.roles

def perform_search(params):
    base_url = "http://localhost:3000"
    search_endpoint = "/search"

    response = requests.get(f"{base_url}{search_endpoint}", params=params)

    if response.status_code == 200:
        logs = response.json()
        for log in logs:
            print(f"Log ID: {log['id']}")
            print(f"Level: {log['level']}")
            print(f"Message: {log['message']}")
            print(f"Resource ID: {log['resourceId']}")
            print(f"Timestamp: {log['timestamp']}")
            print(f"Trace ID: {log['traceId']}")
            print(f"Span ID: {log['spanId']}")
            print(f"Commit: {log['commit']}")
            print(f"Parent Resource ID: {log['parentResource']}")
            print("\n" + "-" * 50 + "\n")
    else:
        print(f"Error: {response.status_code} - {response.text}")

@click.group()
def cli():
    pass

@cli.command()
@click.option('--username', prompt='Enter your username', help='Your username')
@click.password_option(prompt='Enter your password', help='Your password')
def login(username, password):
    user = authenticate(username, password)
    if user:
        click.echo(f"Successfully logged in as {user.username}")
        main_menu(user)
    else:
        click.echo("Authentication failed. Please check your username and password.")

def main_menu(user):
    click.echo("Welcome to the main menu!")
    click.echo(f"You are logged in as {user.username}")
    while True:
        click.echo("\n1. Search logs")
        click.echo("2. Logout")
        choice = click.prompt("Choose an option (1 or 2)", type=int)
        if choice == 1:
            search_menu()
        elif choice == 2:
            click.echo("Logging out...")
            break
        else:
            click.echo("Invalid choice. Please try again.")

def search_menu():
    click.echo("\nSearch logs:")
    level = click.prompt("Enter log level (optional)", default="")
    message = click.prompt("Enter log message (optional)", default="")
    resource_id = click.prompt("Enter resource ID (optional)", default="")
    timestamp = click.prompt("Enter timestamp (optional)", default="")
    trace_id = click.prompt("Enter trace ID (optional)", default="")
    span_id = click.prompt("Enter span ID (optional)", default="")
    commit = click.prompt("Enter commit (optional)", default="")
    parent_resource_id = click.prompt("Enter parent resource ID (optional)", default="")
    start_date = click.prompt("Enter start date (optional, format: YYYY-MM-DD)", default="")
    end_date = click.prompt("Enter end date (optional, format: YYYY-MM-DD)", default="")
    message_regex = click.prompt("Enter message regex (optional)", default="")

    params = {
        'level': level or None,
        'message': message or None,
        'resource_id': resource_id or None,
        'timestamp': timestamp or None,
        'trace_id': trace_id or None,
        'span_id': span_id or None,
        'commit': commit or None,
        'parent_resource_id': parent_resource_id or None,
        'start_date': start_date or None,
        'end_date': end_date or None,
        'message_regex': message_regex or None,
    }

    perform_search(params)

if __name__ == '__main__':
    cli()
