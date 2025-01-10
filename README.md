# Jobarranger unit testing tool using golang

## Register Hosts Command

The `register_hosts` command will get hosts from Zabbix database and categorizing them based on predefined naming Conventions.

### Naming Conventions

The command uses specific naming conventions to determine the type and category of each host:

- **Linux Server**: Hostnames starting with `auto.linux.server.` will be registered as Linux host with the type `server`.
- **Linux Agent**: Hostnames starting with `auto.linux.agent.` will be registered as Linux host with the type `agent`.
- **Windows Agent (not avaliable yet)**: Hostnames starting with `auto.windows.agent.` will be registered as Windows host with the type `agent`.

### Usage

1. **Prerequisites**:
   - Ensure that the Zabbix database is already setup and accessible.

2. **Execution**:
   - Run the command in your terminal or command prompt:
     ```bash
     .\remote_run.exe register_hosts [-p YOUR_DB_HOSTNAME | -m YOUR_DB_HOSTNAME]
     ```
   - `-p`: Specify the hostname of your postgresql database.
   - `-m`: Specify the hostname of your mysql database..
   
   note: Since remote_run doesn't not support multiple database yet, you can choose only one flag.

3. **Verification**:
   - Check the `hosts.json` file in the parent folder. If the registration is successful, the registered hosts will be listed there.
   - If the hosts are registered properly, you will be able to use them from common.Host_pool slice.