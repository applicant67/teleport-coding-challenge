# Teleport Interview Challenge

Submisson for the Teleport L3 Security/ Automation challenge

## Challenge Description/requirements

Description and requirements are copied from:
https://github.com/gravitational/careers/blob/main/challenges/security-automation/challenge.md

### Summary

Build an MVP for workflow automation that can securely configure, add web applications, and users to an Auth0 tenant.

### Rationale

This exercise has two goals:

- It helps us to understand what to expect from you as a Security and Automation Engineer, how you securely configure infrastructure and write code and scripts.
- It helps you get a feel for what it would be like to work at Teleport, as this exercise aims to simulate our day-as-usual and expose you to the type of work we're doing here.

We believe this technique is not only better, but also is more fun compared to whiteboard/quiz interviews so common in the industry. It's not without the downsides - it could take longer than traditional interviews.

Some of the best teams use coding challenges.

We appreciate your time and are looking forward to hack on this project together.

### Design Doc

Before working on the implementation, we ask that you write a brief design document. The design document should cover: scope, infrastructure configuration, APIs you will use, security considerations, edge cases, and implementation details where appropriate.

At Teleport, we prefer Markdown for our designs.

Please submit the design document and all code in a GitHub repository. Public or private is your choice. Please submit the design document as a Pull Request (PR) to allow us to provide you feedback on the proposed design.

A few notes about the design document:

- Try to get the design document approved within the first 2-3 days. This is to ensure you have enough time to work on the implementation.
- Avoid writing an overly detailed design document. One to two pages is sufficient.
- Avoid sending us draft design documents. It is difficult to evaluate which parts are draft and which parts are complete. Instead we encourage asking questions in Slack and sharing a design document that is ready to be reviewed.

Once the design document has been approved by two reviewers, move on to the implementation.

### Implementation

Split your submission into roughly 3-5 Pull Requests to give the team an
opportunity to review your code and provide feedback. Feel free to merge each
PR after you have two approvals.

Our team will do their best to provide a high quality review of the submitted
Pull Requests in a reasonable time frame. You are spending your time on this,
we are going to contribute our time too.

After the final submission, the interview team will assemble and vote using a
"+1, -2" anonymous voting system: +1 is submitted whenever a team member
accepts the submission, -2 otherwise.

In case of a positive result, we will connect you to our HR and recruiting
teams, who will work out the details and present an offer.

In case of a negative score result, the hiring manager will contact you and
share a list of the key observations from the team that affected the result.

### Requirements

Use [GitHub Actions](https://github.com/features/actions) to build MVP workflows
that can do the following.

- Add a workflow that triggers [Terraform](https://www.terraform.io) to
  configure an Auth0 tenant. Enforce strong authentication for all web
  applications within tenant.
- Add a workflow that triggers Terraform to add a web application to the Auth0
  tenant. Use [OSS Grafana](https://grafana.com/oss) with OIDC authentication
  as the sample application.
- Add a workflow that executes a shell script (or Go program) to register
  users within the Auth0 tenant. Add a sample user to the Auth0 tenant.
- Securely store and scope all credentials used.

Link to and share successful workflow runs after each PR has been approved and
merged.
