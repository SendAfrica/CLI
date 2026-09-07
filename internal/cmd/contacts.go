package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

// contactsCmd is the parent command for contact-list management.
var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contact lists, contacts, and Google Contacts sync",
}

var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contact lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/contact-lists"})
		if err != nil {
			return err
		}
		var result []api.ContactList
		if err := decodeListOrPaginated(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new contact list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		req := api.CreateContactListRequest{Name: args[0]}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/contact-lists",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.ContactList
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsUpdateCmd = &cobra.Command{
	Use:   "update [listId]",
	Short: "Update a contact list name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		req := api.UpdateContactListRequest{Name: name}
		data, err := c.Do(client.RequestOpts{
			Method: "PUT",
			Path:   "/v1/contact-lists/" + args[0],
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsDeleteCmd = &cobra.Command{
	Use:   "delete [listId]",
	Short: "Delete a contact list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "DELETE",
			Path:   "/v1/contact-lists/" + args[0],
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsDuplicateCheckCmd = &cobra.Command{
	Use:   "duplicate-check [listId]",
	Short: "Check for duplicate contacts in a list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "GET",
			Path:   "/v1/contact-lists/" + args[0] + "/duplicate-check",
		})
		if err != nil {
			return err
		}
		var result api.DuplicateCheckResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsImportCmd = &cobra.Command{
	Use:   "import [listId] [csv-file]",
	Short: "Import contacts from a CSV file into a contact list",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		listID := args[0]
		csvPath := args[1]

		fileContent, err := os.ReadFile(csvPath)
		if err != nil {
			return fmt.Errorf("reading CSV file: %w", err)
		}

		body, contentType, err := client.MultipartForm(
			map[string]string{},
			"file",
			filepathBase(csvPath),
			fileContent,
			"text/csv",
		)
		if err != nil {
			return err
		}

		c, err := apiClient()
		if err != nil {
			return err
		}

		opts := client.RequestOpts{
			Method: "POST",
			Path:   "/v1/contact-lists/" + listID + "/import",
			Body: &csvBody{
				buffer:      body,
				contentType: contentType,
			},
		}
		data, err := c.Do(opts)
		if err != nil {
			return err
		}
		var result api.ContactImportResult
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

// csvBody adapts *bytes.Buffer to the client's body encoding.
type csvBody struct {
	buffer      *bytes.Buffer
	contentType string
}

func (b *csvBody) Bytes() []byte { return b.buffer.Bytes() }

func (b *csvBody) ContentType() string { return b.contentType }

func filepathBase(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[i+1:]
		}
	}
	return p
}

var contactsListContactsCmd = &cobra.Command{
	Use:   "list-contacts [listId]",
	Short: "List contacts in a contact list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		params := map[string]string{}
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		if page > 0 {
			params["page"] = fmt.Sprintf("%d", page)
		}
		if perPage > 0 {
			params["per_page"] = fmt.Sprintf("%d", perPage)
		}
		data, err := c.Do(client.RequestOpts{
			Method:      "GET",
			Path:        "/v1/contact-lists/" + args[0] + "/contacts",
			QueryParams: params,
		})
		if err != nil {
			return err
		}
		var result []api.Contact
		if _, err := decodePaginated(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsCreateContactCmd = &cobra.Command{
	Use:   "add-contact [listId] [phone1,phone2,...]",
	Short: "Add a contact to a contact list",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		listID := args[0]
		phones := parseStringSlice(args[1])
		name, _ := cmd.Flags().GetString("name")

		req := api.CreateContactRequest{
			Name:   name,
			Phones: phones,
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/contact-lists/" + listID + "/contacts",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.Contact
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsGetContactCmd = &cobra.Command{
	Use:   "get-contact [listId] [contactId]",
	Short: "Get a single contact by ID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "GET",
			Path:   "/v1/contact-lists/" + args[0] + "/contacts/" + args[1],
		})
		if err != nil {
			return err
		}
		var result api.Contact
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsUpdateContactCmd = &cobra.Command{
	Use:   "update-contact [listId] [contactId]",
	Short: "Update a contact",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		phonesStr, _ := cmd.Flags().GetString("phones")

		req := api.UpdateContactRequest{}
		if name != "" {
			req.Name = name
		}
		if phonesStr != "" {
			req.Phones = parseStringSlice(phonesStr)
		}

		data, err := c.Do(client.RequestOpts{
			Method: "PUT",
			Path:   "/v1/contact-lists/" + args[0] + "/contacts/" + args[1],
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsDeleteContactCmd = &cobra.Command{
	Use:   "delete-contact [listId] [contactId]",
	Short: "Delete a contact",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "DELETE",
			Path:   "/v1/contact-lists/" + args[0] + "/contacts/" + args[1],
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsExportCmd = &cobra.Command{
	Use:   "export [listId]",
	Short: "Export contacts from a list as CSV",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		resp, env, err := c.DoWithResponse(client.RequestOpts{
			Method: "GET",
			Path:   "/v1/contact-lists/" + args[0] + "/contacts/export",
		})
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if env.Error != nil {
			return env.Error
		}

		if printer().Format() == "table" {
			fmt.Fprint(os.Stderr, "Use --output json or --output xml to export CSV-friendly data.\n")
		}
		return printer().Print(rawMessage(env.Data))
	},
}

var contactsAddPhoneCmd = &cobra.Command{
	Use:   "add-phone [listId] [contactId] [phone]",
	Short: "Add a phone number to a contact",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		label, _ := cmd.Flags().GetString("label")
		req := api.AddContactPhoneRequest{
			Phone: args[2],
			Label: label,
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/contact-lists/" + args[0] + "/contacts/" + args[1] + "/phones",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsDeletePhoneCmd = &cobra.Command{
	Use:   "delete-phone [listId] [contactId] [phoneId]",
	Short: "Delete a phone number from a contact",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "DELETE",
			Path:   "/v1/contact-lists/" + args[0] + "/contacts/" + args[1] + "/phones/" + args[2],
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsGoogleStatusCmd = &cobra.Command{
	Use:   "google-status",
	Short: "Check Google Contacts sync connection status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "GET",
			Path:   "/v1/contact-lists/google/status",
		})
		if err != nil {
			return err
		}
		var result api.GoogleContactsStatus
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var contactsGoogleSyncCmd = &cobra.Command{
	Use:   "google-sync",
	Short: "Sync contacts from Google Contacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/contact-lists/google/sync",
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var contactsGoogleDisconnectCmd = &cobra.Command{
	Use:   "google-disconnect",
	Short: "Disconnect Google Contacts sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "DELETE",
			Path:   "/v1/contact-lists/google/disconnect",
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

func init() {
	contactsCreateContactCmd.Flags().String("name", "", "contact name")
	contactsUpdateContactCmd.Flags().String("name", "", "new contact name")
	contactsUpdateContactCmd.Flags().String("phones", "", "comma-separated phone numbers to set")
	contactsAddPhoneCmd.Flags().String("label", "", "phone label: mobile/home/work/other")
	contactsListContactsCmd.Flags().Int("page", 1, "page number")
	contactsListContactsCmd.Flags().Int("per-page", 50, "items per page")
}
