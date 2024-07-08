<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <div>
    <div class="container">
      <section class="section">
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li><router-link :to="{ name: 'Apps' }">Apps</router-link></li>
            <li class="is-active">
              <router-link :to="{ name: 'Costs' }">Costs</router-link>
            </li>
          </ul>
          <b-button
            @click="CopyHrefToClipboard()"
            icon-left="content-copy"
            size="is-small"
          ></b-button>
        </nav>
        <h1 class="title is-1">Costs</h1>
        <p class="subtitle is-3">Keep track of spending</p>
        <div v-if="notes !== '' || canUserAccountAdmin">
          <div class="content">
            <label class="label">Notes</label>
            <p
              :class="
                canUserAccountAdmin
                  ? 'display-is-editable pointer-cursor-on-hover'
                  : ''
              "
              class="subtitle is-4 notes-highlight"
              @click="EditCostsNotes"
            >
              <i>
                {{ notes || "Add notes" }}
              </i>
            </p>
          </div>
          <br />
        </div>
        <div>
          <h1>Report</h1>
          <div
            id="cost-view-chart-frame"
            style="position: relative; height: 100%; width: 100%"
          >
            <canvas id="costs-view-chart"></canvas>
          </div>
        </div>
        <div class="columns">
          <div class="card column">
            <div class="media">
              <div class="media-left">
                <p class="subtitle">daily cost average</p>
              </div>
              <div class="media-content">
                <p class="subtitle">${{ view.totalDailyCostAverage }}</p>
              </div>
            </div>
          </div>
          <div class="card column">
            <div class="media">
              <div class="media-left">
                <p class="subtitle">weekly cost average</p>
              </div>
              <div class="media-content">
                <p class="subtitle">${{ view.totalWeeklyCostAverage }}</p>
              </div>
            </div>
          </div>
          <div class="card column">
            <div class="media">
              <div class="media-left">
                <p class="subtitle">three-monthly average</p>
              </div>
              <div class="media-content">
                <p class="subtitle">${{ view.totalThreeMonthAverage }}</p>
              </div>
            </div>
          </div>
          <div class="card column">
            <div class="media">
              <div class="media-left">
                <p class="subtitle">yearly spend</p>
              </div>
              <div class="media-content">
                <p class="subtitle">${{ view.totalYearCumulativeSpend }}</p>
              </div>
            </div>
          </div>
        </div>
        <br />
        <div>
          <label class="label">Search for lists</label>
          <b-field>
            <b-input
              icon="magnify"
              size="is-medium"
              placeholder="Enter a list name"
              type="search"
              expanded
              v-model="listSearch"
              ref="search"
            >
            </b-input>
            <p class="control">
              <b-select
                placeholder="Sort by"
                icon="sort"
                v-model="sortBy"
                size="is-medium"
                expanded
              >
                <option value="recentlyAdded">Recently Added</option>
                <option value="lastAdded">Last Added</option>
                <option value="recentlyUpdated">Recently Updated</option>
                <option value="lastUpdated">Last Updated</option>
                <option value="alphabeticalDescending">A-z</option>
                <option value="alphabeticalAscending">z-A</option>
              </b-select>
            </p>
          </b-field>
          <b-loading
            :is-full-page="false"
            :active.sync="pageLoading"
            :can-cancel="false"
          ></b-loading>
          <section>
            <div
              class="card pointer-cursor-on-hover"
              @click="
                $router.push({
                  name: 'New cost entry',
                })
              "
            >
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon icon="cart-plus" size="is-medium"> </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">Add a new cost entry</p>
                  </div>
                  <div class="media-right">
                    <b-icon
                      icon="chevron-right"
                      size="is-medium"
                      type="is-midgray"
                    ></b-icon>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <floatingAddButton
          :routerLink="{
            name: 'New cost entry',
          }"
          v-if="displayFloatingAddButton"
        />
        <br />
        <div v-if="viewItemsDisplay.length > 0">
          <div>
            <b-table
              :data="viewItemsDisplay"
              :columns="viewTableColumns"
              focusable
              striped
              paginated
              per-page="20"
            >
            </b-table>
          </div>
          <br />
          <p>{{ viewItemsDisplay.length }} entries displayed</p>
        </div>
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="cart-remove"
                    size="is-medium"
                    type="is-midgray"
                  ></b-icon>
                </div>
                <div class="media-content">
                  <p
                    class="subtitle is-4"
                    v-if="
                      listSearch === '' && lists.length === 0 && !pageLoading
                    "
                  >
                    No lists added yet.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="
                      listSearch === '' &&
                      listDisplayState === 1 &&
                      lists.length > 0 &&
                      !pageLoading
                    "
                  >
                    All lists have been completed.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="
                      listSearch === '' &&
                      listDisplayState === 2 &&
                      lists.length > 0 &&
                      !pageLoading
                    "
                  >
                    No lists have been completed yet.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="listSearch !== '' && !pageLoading"
                  >
                    No lists found.
                  </p>
                  <p class="subtitle is-4" v-else-if="pageLoading">
                    Loading lists...
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from "@/common/common";
import costs from "@/requests/authenticated/costs";
import flatmates from "@/requests/authenticated/flatmates";
import cani from "@/requests/authenticated/can-i";
import { DialogProgrammatic as Dialog } from "buefy";
import { Chart } from "chart.js/auto";

export default {
  name: "costs-view",
  data() {
    return {
      viewTableColumns: [
        {
          field: "title",
          label: "Title",
        },
        {
          field: "paymentType",
          label: "Type",
        },
        {
          field: "amount",
          label: "Amount",
        },
        {
          field: "invoiceDate",
          label: "Invoice Date",
        },
      ],
      displayFloatingAddButton: true,
      canUserAccountAdmin: false,
      notes: "",
      users: {},
      view: {},
      chart: {},
      chartData: {
        type: "line",
        data: {
          labels: [
            "three months ago",
            "two months ago",
            "last month",
            "this month",
          ],
          datasets: [
            {
              label: "Spending over the last three months",
              data: [],
            },
          ],
        },
        options: {
          responsive: true,
          lineTension: 0.4,
          scales: {
            y: {
              ticks: {
                callback: (value, index, ticks) => "$" + value,
                stepSize: 10,
              },
            },
          },
        },
      },
      deviceIsMobile: false,
      pageLoading: true,
      sortBy: "recentlyUpdated",
    };
  },
  components: {
    floatingAddButton: () =>
      import("@/components/common/floating-add-button.vue"),
  },
  computed: {
    viewItemsDisplay() {
      // TODO express invoicedDate, invoicedBy, author, authorLast
      // TODO use this.users as id-to-names cache
      return (
        this.viewItemsFiltered.map((i) => {
          if (i.invoiceDate !== 0) {
            i.invoiceDate = common.TimestampToCalendar(i.invoiceDate);
          }
          i.amount = common.FormatFloatAsMoney(i.amount);
          return i;
        }) || []
      );
    },
    viewItemsFiltered() {
      return this.view.items || [];
    },
  },
  methods: {
    CopyHrefToClipboard() {
      common.CopyHrefToClipboard();
    },
    GetView() {
      costs
        .GetView()
        .then((resp) => {
          this.pageLoading = false;
          this.view = resp.data.data || {};
        })
        .catch(() => {
          common.DisplayFailureToast(
            "Hmmm seems somethings gone wrong loading the costs view"
          );
        });
    },
    GetCostsNotes() {
      costs
        .GetCostsNotes()
        .then((resp) => {
          this.notes = resp.data.spec || "";
          this.pageLoading = false;
        })
        .catch(() => {
          common.DisplayFailureToast(
            "Hmmm seems somethings gone wrong loading the notes for costs"
          );
        });
    },
    EditCostsNotes() {
      if (this.canUserAccountAdmin !== true) {
        return;
      }
      this.displayFloatingAddButton = false;
      Dialog.prompt({
        title: "Cost notes",
        message: `Enter notes that are useful for keeping track of costs.`,
        container: null,
        icon: "text",
        hasIcon: true,
        inputAttrs: {
          placeholder: "e.g. make sure to add rent costs.",
          maxlength: 250,
          required: false,
          value: this.notes || undefined,
        },
        trapFocus: true,
        onConfirm: (value) => {
          // shoppinglist
          //   .PutShoppingListNotes(value)
          //   .then(() => {
          //     this.pageLoading = true;
          //     this.displayFloatingAddButton = true;
          //     this.GetShoppingListNotes();
          //   })
          //   .catch((err) => {
          //     common.DisplayFailureToast(
          //       "Failed to update notes" +
          //         `<br/>${err.response.data.metadata.response}`
          //     );
          //     this.displayFloatingAddButton = true;
          //   });
        },
        onCancel: () => {
          this.displayFloatingAddButton = true;
        },
      });
    },
    GetFlatmateName(id) {
      flatmates
        .GetFlatmate(id)
        .then((resp) => {
          return resp.data.spec.names;
        })
        .catch((err) => {
          common.DisplayFailureToast(
            "Failed to fetch user account" +
              `<br/>${err.response.data.metadata.response}`
          );
          return id;
        });
    },
    CheckDeviceIsMobile() {
      this.deviceIsMobile = common.DeviceIsMobile();
    },
  },
  watch: {
    view() {
      if (this.view === null) {
        return;
      }
      const ctx = document.getElementById("costs-view-chart");
      this.chartData.data.datasets[0].data = [
        this.view.threeMonthsAgoCumulative,
        this.view.twoMonthsAgoCumulative,
        this.view.oneMonthAgoCumulative,
        this.view.thisMonthCumulative,
      ];
      this.chart = new Chart(ctx, this.chartData);
    },
  },
  async beforeMount() {
    cani.GetCanIgroup("admin").then((resp) => {
      this.canUserAccountAdmin = resp.data.data;
    });
    this.GetView();
  },
  mounted() {},
  beforeDestroy() {
    window.removeEventListener("resize", this.CheckDeviceIsMobile, true);
  },
  async created() {
    this.CheckDeviceIsMobile();
    window.addEventListener("resize", this.CheckDeviceIsMobile, true);
  },
};
</script>

<style scoped></style>
