Site Members
============

Allows site admins to grant or deny access to a site based on the email address of the user.

Just edit the config file so that the following values are set:

`Subdomain` = the subdomain for a site, this is as per your Microcosm control panel, if your site was originally https://foo.microco.sm/ then your subdomain is `foo`. If you have a custom domain name, such as www.foobar.com, your subdomain is **still** `foo`.

`Token` = an `access_token` for a privileged user on the site, either an administrator, site owner or moderator. If you are one of these people then you can view your cookie using Google Chrome when signed into the site by pressing F12 and then choosing Resources > Cookies. The access_token is a very long string.

`IsMember` = a true|false value, if true it adds the person as a member, if false it removes them as a member

`Emails` = a comma-delimited, "quoted" list of email addresses to add or remove

An example:

```
{
	"Subdomain": "foo",
	"Token": "access_token",
	"IsMember": true,
	"Emails": [
		"someone@example.com",
		"someone-else@example.com"
	]
}
```

Then just:

`go build && ./site-members`
