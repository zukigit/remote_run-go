# Jobarranger unit testing tool framework using golang

## register_hosts Command

The `register_hosts` command will get hosts from Zabbix database and categorizing them based on predefined naming Conventions.

### Naming Conventions

The command uses specific naming conventions to determine the type and category of each host:

- **Linux Server**: Hostnames starting with `auto.linux.server.` will be registered as Linux host with the type `server`.
- **Linux Agent**: Hostnames starting with `auto.linux.agent.` will be registered as Linux host with the type `agent`.
- **Windows Agent (not avaliable yet)**: Hostnames starting with `auto.windows.agent.` will be registered as Windows host with the type `agent`.

### Usage

1. **Prerequisites**:
   - Ensure that the Zabbix database is already setup and accessible.
   - If the host is windows, check **PubkeyAuthentication** in `C:\ProgramData\ssh\sshd_config`. If it is no, set it to yes and restart **ssh-server**.

2. **Execution**:
   - Run the command in your terminal or command prompt:
     ```bash
     .\remote_run.exe register_hosts [-p YOUR_DB_HOSTNAME | -m YOUR_DB_HOSTNAME]
     ```
   - `-p`: Specify the hostname of your postgresql database.
   - `-m`: Specify the hostname of your mysql database.
   - `--db-user`: Specify the database username to connect. (OPTIONAL)
   - `--db-password`: Specify the database password to connect. (OPTIONAL)
   - `--db-name`: Specify the database name to connect. (OPTIONAL)
   - `--db-port`: Specify the database port to connect. Default: **5432(psql) | 3306(mysql)**. (OPTIONAL)

   note: Since **remote_run** doesn't not support multiple databases yet, you can choose only one flags between **-p** and **-m**.

3. **Verification**:
   - Check the `hosts.json` file in the parent folderpath. If the registration is successful, the registered hosts will be listed there.
   - If the hosts are registered properly, you will be able to use them from **common.Hosts** slice.

## remote_run Command

The `remote_run` command will run specific ticket.

### Usage

1. **Prerequisites**:
   - Ensure target hosts are registered in **hosts.json** file.

2. **Execution**:
   - Run the command in your terminal or command prompt:
     ```bash
     .\remote_run.exe TICKET_NUMBER [-p YOUR_DB_HOSTNAME | -m YOUR_DB_HOSTNAME]
     ```
   - `-p`: Specify the hostname of your postgresql database.
   - `-m`: Specify the hostname of your mysql database.
   - `-a`: Use this flag to run all avaliable tickets. Default: **false**. (OPTIONAL)
   - `--testcase`: Specify testcase number to run specific testcase. Default: **0**. (OPTIONAL)
   - `--db-user`: Specify the database username to connect. (OPTIONAL)
   - `--db-password`: Specify the database password to connect. (OPTIONAL)
   - `--db-name`: Specify the database name to connect. (OPTIONAL)
   - `--db-port`: Specify the database port to connect. Default: **5432(psql) | 3306(mysql)**. (OPTIONAL)
   - `--timeout`: Specify common timout in seconds. Default: **300**. (OPTIONAL)

   note: Since **remote_run** doesn't not support multiple databases yet, you can choose only one flags between **-p** and **-m**.
