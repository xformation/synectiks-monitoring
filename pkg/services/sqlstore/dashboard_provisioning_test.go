package sqlstore

import (
	"testing"
	"time"

	"github.com/xformation/synectiks-monitoring/pkg/components/simplejson"
	"github.com/xformation/synectiks-monitoring/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDashboardProvisioningTest(t *testing.T) {
	Convey("Testing Dashboard provisioning", t, func() {
		InitTestDB(t)

		folderCmd := &models.SaveDashboardCommand{
			OrgId:    1,
			FolderId: 0,
			IsFolder: true,
			Dashboard: simplejson.NewFromAny(map[string]interface{}{
				"id":    nil,
				"title": "test dashboard",
			}),
		}

		err := SaveDashboard(folderCmd)
		So(err, ShouldBeNil)

		saveDashboardCmd := &models.SaveDashboardCommand{
			OrgId:    1,
			IsFolder: false,
			FolderId: folderCmd.Result.Id,
			Dashboard: simplejson.NewFromAny(map[string]interface{}{
				"id":    nil,
				"title": "test dashboard",
			}),
		}

		Convey("Saving dashboards with provisioning meta data", func() {
			now := time.Now()

			cmd := &models.SaveProvisionedDashboardCommand{
				DashboardCmd: saveDashboardCmd,
				DashboardProvisioning: &models.DashboardProvisioning{
					Name:       "default",
					ExternalId: "/var/grafana.json",
					Updated:    now.Unix(),
				},
			}

			err := SaveProvisionedDashboard(cmd)
			So(err, ShouldBeNil)
			So(cmd.Result, ShouldNotBeNil)
			So(cmd.Result.Id, ShouldNotEqual, 0)
			dashId := cmd.Result.Id

			Convey("Can query for provisioned dashboards", func() {
				query := &models.GetProvisionedDashboardDataQuery{Name: "default"}
				err := GetProvisionedDashboardDataQuery(query)
				So(err, ShouldBeNil)

				So(len(query.Result), ShouldEqual, 1)
				So(query.Result[0].DashboardId, ShouldEqual, dashId)
				So(query.Result[0].Updated, ShouldEqual, now.Unix())
			})

			Convey("Can query for one provisioned dashboard", func() {
				query := &models.IsDashboardProvisionedQuery{DashboardId: cmd.Result.Id}

				err := GetProvisionedDataByDashboardId(query)
				So(err, ShouldBeNil)

				So(query.Result, ShouldBeTrue)
			})

			Convey("Can query for none provisioned dashboard", func() {
				query := &models.IsDashboardProvisionedQuery{DashboardId: 3000}

				err := GetProvisionedDataByDashboardId(query)
				So(err, ShouldBeNil)
				So(query.Result, ShouldBeFalse)
			})

			Convey("Deleteing folder should delete provision meta data", func() {
				deleteCmd := &models.DeleteDashboardCommand{
					Id:    folderCmd.Result.Id,
					OrgId: 1,
				}

				So(DeleteDashboard(deleteCmd), ShouldBeNil)

				query := &models.IsDashboardProvisionedQuery{DashboardId: cmd.Result.Id}

				err = GetProvisionedDataByDashboardId(query)
				So(err, ShouldBeNil)
				So(query.Result, ShouldBeFalse)
			})
		})
	})
}
