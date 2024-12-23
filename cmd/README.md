# Commands

## Register Hosts Command

The `register_hosts` command will get hosts from Zabbix database and categorizing them based on predefined naming Conventions.

### Naming Conventions

The command uses specific naming conventions to determine the type and category of each host:

- **Linux Server**: Hostnames starting with `auto.linux.server.` are registered as Linux hosts with the type `server`.
- **Linux Agent**: Hostnames starting with `auto.linux.agent.` are registered as Linux hosts with the type `agent`.
- **Windows Agent**: Hostnames starting with `auto.windows.agent.` are registered as Windows hosts with the type `agent`.

### Usage

1. **Prerequisites**:
   - Ensure that the Zabbix database is already setup and accessible.

2. **Execution**:
   - Run the command in your terminal or command prompt:
     ```bash
     .\remote_run.exe register_hosts --db-hostname YOUR_DB_HOSTNAME [--with-postgresql | --with-mysql]
     ```
   - `--db-hostname`: Specify the hostname of your database. This is mandatory.
   - `--with-postgresql`: Use this flag if you're working with a PostgreSQL database.
   - `--with-mysql`: Use this flag if you're working with a MySQL database.

3. **Verification**:
   - Check the `hosts.json` file in the parent folder. If the registration is successful, the registered hosts will be listed there.