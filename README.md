# Queries

This repository serves as a storage and management hub for **CSIEM (Cognitive SIEM)** and **SPL (Splunk Processing Language)** queries. The goal is to provide a well-organized collection of queries that can be used to enhance security operations and incident detection in both **CSIEM** and **Splunk** environments.

## Table of Contents

- [What’s New](#whats-new)
- [Features](#features)
- [How to Use](#how-to-use)
- [Acknowledgements](#acknowledgements)
- [Contributing](#contributing)
- [License](#license)

## What’s New

Currently, this repository includes Sigma rules converted into **CSIEM** and **SPL** query languages. The conversion is done using [Alterix](https://github.com/mtnmunuklu/alterix), a tool designed to translate Sigma rules into formats compatible with both **CSIEM** and **Splunk**. This transformation aims to streamline the integration of Sigma rules into both platforms, enabling efficient querying and analysis.

## Features

- **Sigma to CSIEM Conversion**: Sigma rules have been successfully converted to **CSIEM** query language.
- **Sigma to SPL Conversion**: Sigma rules have been converted into **SPL** for use within **Splunk** environments.
- **Repository Organization**: Queries are organized by their respective query languages for easy management and usage.
- **Future Updates**: More features and query enhancements will be added as the project evolves.

## How to Use

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/mtnmunuklu/queries.git

2. **Browse Queries**: All `queries` are available in the queries directory. The queries are organized by the language type (CSIEM or SPL), allowing users to easily find the desired query format.

3. **Integration**:
   - For **CSIEM**: Integrate the CSIEM queries into your **CSIEM** platform to improve your organization’s security monitoring and threat detection capabilities.
   - For **Splunk**: Use the **SPL** queries in your **Splunk** environment to enhance your log analysis and security operations.

## Contributing

Contributions to Queries are welcome and encouraged! Please read the [contribution guidelines](CONTRIBUTING.md) before making any contributions to the project.

## Acknowledgements

We would like to extend our gratitude to the creators and contributors of [Sigma](https://github.com/Neo23x0/sigma) for their work in developing a powerful and flexible format for representing detection rules. 

Additionally, we have included the names of the original authors of all converted Sigma rules in the `description` field of the corresponding  CSIEM and SPL queries, ensuring proper credit is given for their contributions to the community.

## License

Queries is licensed under the MIT License. See [LICENSE](LICENSE) for the full text of the license.