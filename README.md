# Terraform Google Calendar Provider

This is a [Terraform][terraform] provider for managing meetings on Google
Calendar. It enables you to treat "calendars as code" the same way you already
treat infrastructure as code!


## Installation

1. Download the latest compiled binary from [GitHub releases][releases].

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins/linux_amd64
    $ mv terraform-provider-googlecalendar $HOME/.terraform.d/plugins/terraform-provider-googlecalendar_v0.1.0
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```

    This will find the plugin locally.


## Usage

1. You will need a valid Google Cloud account and permission to create a service
account. You can create the service account in any project, but make sure you
choose "server-to-server" communication, since this is not an OAuth application.

    1. Visit the [Google Cloud Credentials Console][gcloud-creds]
    1. Click "Create credentials"
    1. Choose "Service account key"
    1. Use "Compute Engine default service account" (or make your own)
    1. Choose "JSON" as the key type
    1. Click "Create"

    After a few seconds, your browser will download a credentials file in JSON.
    Save this file securely (treat it like a password).

    Note: if you are making your own service account, be sure to grant access to
    the "calendar" OAuth scope!

1. Set the `CALENDAR_CREDENTIALS_FILE` environment variable to point to your newly-downloaded credentials file:

    ```sh
    $ export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
    ```

    The Terraform provider automatically reads this environment variable and
    uses the file at the given path for authentication.

1. Set the `CALENDAR_OAUTH2_TOKEN_FILE` environment variable to point to your personal token file obtained following [this](https://developers.google.com/calendar/quickstart/go) instructions.


1. Create a Terraform configuration file:

    ```hcl
    resource "googlecalendar_event" "example" {

      // The calendar_id is optional, it can be used to manage other calendar
      // than the user primary. Use it when managing a shared calendar.
      calendar_id = "primary"

      summary     = "My Event"
      description = "Long-form description of the event"
      location    = "Conference Room B"

      // Start and end times work best if specified as RFC3339.
      // To specify a full day event, use "2017-10-12" format.
      start = "2017-10-12T15:00:00-05:00"
      end   = "2017-10-12T17:00:00-05:00"

      // Each attendee is listed separately, and attendees can be marked as
      // optional.
      attendee {
        email = "seth@sethvargo.com"
      }

      attendee {
        email    = "you@company.com"
        optional = true
      }
    }
    ```

1. Run `terraform init` to pull in the provider:

    ```sh
    $ terraform init
    ```

1. Run `terraform plan` and `terraform apply` to create events:

    ```sh
    $ terraform plan

    $ terraform apply
    ```

## Examples

For more examples, please see the [examples][examples] folder in this
repository.

## Reference

### Arguments

Arguments are provided as inputs to the resource, in the `*.tf` file.

- `summary` `(string, required)` - the "title" of the event.

- `calendar_id` `(string, optiona)` - the calendar id in which to manage the
  events. Optional, _primary_ by default.

- `start` `(string, required)` - the RFC3339-formatted start time of the event
  _with the timestamp included_ if it is a DateTime, only the date part to
  define a full day event.

- `end` `(string, required)` - the RFC3339-formatted end time of the event _with
  the timestamp included_ if it is a DateTime, only the date part to
  define a full day event.

- `description` `(string)` - the long-form description of the event. This can be
  multiple paragraphs using Terraform's heredoc syntax.

- `guests_can_invite_others` `(bool, true)` - specifies that guests (attendees)
  can invite other guests. Set this to false to allow only the organizer to
  manage the guest list.

- `guests_can_modify` `(bool, false)` - specifies that guests (attendees) can
  modify the event (change start time, description, etc). Set this to true to
  give any guest full control over the event.

- `guests_can_see_other_guests` `(bool, true)` - specifies that guests
  (attendees) can see other guests. Set this to false to restrict the guest list
  visibility.

- `show_as_available` `(bool, false)` - specifies that the time should be
  "blocked" on the calendar (mark as busy). Set this to true to create an event
  that is transparent.

- `send_notifications` `(bool, true)` - specifies that email notifications
  should be sent to guests (attendees). Set this to false to put things on
  people's calendar's without notifying them.

- `visibility` `(string)` - specifies the visibility for the event. Valid values
  are:

    - `""` - default inherit from calendar
    - `"public"` - public
    - `"private"` - private

- `attendee` `(list of structures)` - specifies a guest (attendee) to invite to
  the event. This may be specified more than once to invite multiple people to
  the same event. The following fields are supported:

    - `email` `(string, required)` - the Google email address of the attendee.

    - `optional` `(bool, false)` - specifies that the guest (attendee) is marked
      as optional. Set this to true to mark the user as an optional attendee.

- `reminder` `(list of reminders)` - specifies a reminder option leading up to
  the event for all attendees. This overrides any default reminders the user has
  set for their calendar. Leave this unset to inherit calendar default
  reminders. This may be specified more than once to remind multiple times. The
  following fields are supported:

    - `method` `(string, required)` - the method to use. Valid options are:

      - `"email"` - send an email
      - `"popup"` - popup in-browser
      - `"sms"` - send a text message (requires GSuite)

    - `before` `(string, required)` - the duration prior to the meeting to send
      the reminder. Note that the Google Calendar API expects this to be the
      "number of minutes", but Terraform adds syntactic sugar here by allowing
      you to specify the time in a Go timestamp like "30m" or "4h". These
      timestamps are parsed and converted to minutes automatically.

### Attributes

Attributes are values that are only known after creation.

- `event_id` `(string)` - the unique ID of the event on this calendar

- `hangout_link` `(string)` - the HTTPS web link to the attached Google Hangout.
  In practice, I have been unable to get this link to appear.

- `html_link` `(string)` - the HTTP web link to the calendar invite on
  [calendar.google.com](https://calendar.google.com/).

## Constraints & Understanding

[terraform]: https://www.terraform.io/
[releases]: https://github.com/sethvargo/terraform-provider-googlecalendar/releases
[gcloud-creds]: https://console.cloud.google.com/apis/credentials
[examples]: https://github.com/sethvargo/terraform-provider-googlecalendar/tree/master/examples
