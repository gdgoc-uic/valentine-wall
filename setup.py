import shutil
import os

setup_env_done = False
default_values = {
    'ENV': 'production',
    'BASE_DOMAIN': 'valentine-wall.gdscuic.org',
    'FRONTEND_URL': 'https://valentine-wall.gdscuic.org',
    'BACKEND_URL': 'https://valentine-wall.gdscuic.org/pb',
}

try:
    print("Setting up .env file...")

    # check if .env file does not exist
    if not os.path.exists('.env'):
        # copy .env.example to .env
        shutil.copy('.env.example', '.env')
    else:
        setup_env_done = True
        print(".env file already exists.")

    # open .env file
    with open('.env', 'r') as file:
        # read a list of lines into data
        data = file.readlines()

        # uncomment and change the value of the commented environment variables
        for idx, line in enumerate(data):
            if line.startswith('#') and line.endswith('=\n'):
                # remove the comment
                data[idx] = line[2:]
                line = data[idx]

                # get name of the environment variable
                env_var_name = line.split('=')[0]

                while True:
                    # prompt value of the environment variable
                    env_var_value = input(f'Enter value for {env_var_name}: ')

                    if len(env_var_value) == 0:
                        # if default value is empty and user is empty, prompt again
                        print(f'Value for {env_var_name} is required.')
                        continue

                    break

                # change the value of the environment variable
                data[idx] = f'{env_var_name}={env_var_value}\n'
            elif (line.startswith('#') and not line.endswith('=\n')) or len(line.strip()) == 0:
                continue
            else:
                kv = line.split('=')

                # check if the environment variable is not commented
                env_var_name = kv[0]
                env_var_value = kv[1].strip()

                if env_var_name in default_values and env_var_value != default_values[env_var_name]:
                    env_var_value = default_values[env_var_name]

                while True:
                    # prompt value of the environment variable
                    if len(env_var_value) == 0:
                        env_var_value = input(f'Enter value for {env_var_name}: ')

                        print(len(env_var_value))
                    else:
                        env_var_value2 = input(f'Enter value for {env_var_name} (Default: {env_var_value}): ')
                        if len(env_var_value2) != 0:
                            env_var_value = env_var_value2

                    if len(env_var_value) == 0:
                        # if default value is empty and user is empty, prompt again
                        print(f'Value for {env_var_name} is required.')
                        continue

                    break

                # change the value of the environment variable
                data[idx] = f'{env_var_name}={env_var_value}\n'

    # write data to .env file
    with open('.env', 'w') as file:
        file.writelines(data)
        setup_env_done = True

    print(".env file setup done.")

    # Execute ./deploy.sh
    print("Deploying...")
    os.system('./deploy.sh')
except KeyboardInterrupt:
    if not setup_env_done:
        # Remove .env file
        os.remove('.env')
    print("Setup cancelled.")
