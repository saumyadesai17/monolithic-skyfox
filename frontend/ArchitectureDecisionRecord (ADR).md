# Architecture Decision Record (ADR)

ID | Title | Context      | Consequences  | Status  | Link to ADR | 
---| ------------ | ------------ | ------------- | ------- | ----------- | 
ADR-001 | Using React |React capabilities will help in building both webapp and can be extended to be used in mobile interfaces. Skyfox also later on plan to hire React developers | There is a learning curve if developers are new to React | Accepted | |
ADR-002 | Using separate repository for front-end and backend | Skyfox project front-end code built using React which is different from backend (Spring Boot), so having different repo for separate the concerns | We will maintain two repos | Accepted | |
ADR-003 | Use VSCode IDE for front-end development | Have consistent among the developers so that it will help in pairing. VS code is one of best IDE for front-end development |  | Accepted | |
ADR-004 | Use of gitlab | We want to use GitLab as central repo to maintain the code remotely, Git Lab provides free pipeline minutes | What if we ran out of free pipeline minutes? | Deprecated | |
ADR-005 | Use of AWS code commit | We do not want to manage AWS keys in Gitlab, we are exposing admin credentials within Gitlab | We need authenticate developers using OKTA | Superseded | ADR-004 |
ADR-006 | Use saml2aws authentication | To authenticate developers using ThoughtWorks credentials  | Everyone who wants to access the code commit repo needs to configure before login | Accepted | |
