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
      <b-modal
        scroll="keep"
        v-model="isNewItemModalActive"
        :fullscreen="deviceIsMobile"
        has-modal-card
        :can-cancel="false"
      >
        <costsItemNew></costsItemNew>
      </b-modal>
      <b-modal
        scroll="keep"
        v-model="isEditItemModalActive"
        :fullscreen="deviceIsMobile"
        has-modal-card
        :can-cancel="false"
      >
        <costsItemEdit v-bind="editItemProps"></costsItemEdit>
      </b-modal>
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
        <div v-if="notes !== '' || canUserAccountAdmin" class="mb-3">
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
              <i> {{ notes || "Add notes" }} </i>
            </p>
          </div>
        </div>
        <div v-if="view.items !== null">
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
            <div class="card column" v-if="view.totalDailyCostAverage">
              <div class="media">
                <div class="media-content">
                  <p class="subtitle">daily cost average</p>
                  <p class="subtitle">
                    ${{ FmtAsCurrency(view.totalDailyCostAverage) }}
                  </p>
                </div>
              </div>
            </div>
            <div class="card column" v-if="view.totalWeeklyCostAverage">
              <div class="media">
                <div class="media-content">
                  <p class="subtitle">weekly cost average</p>
                  <p class="subtitle">
                    ${{ FmtAsCurrency(view.totalWeeklyCostAverage) }}
                  </p>
                </div>
              </div>
            </div>
            <div class="card column" v-if="view.totalThreeMonthAverage">
              <div class="media">
                <div class="media-content">
                  <p class="subtitle">three-monthly average</p>
                  <p class="subtitle">
                    ${{ FmtAsCurrency(view.totalThreeMonthAverage) }}
                  </p>
                </div>
              </div>
            </div>
            <div class="card column" v-if="view.totalYearCumulativeSpend">
              <div class="media">
                <div class="media-content">
                  <p class="subtitle">yearly spend</p>
                  <p class="subtitle">
                    ${{ FmtAsCurrency(view.totalYearCumulativeSpend) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <floatingAddButton
          v-if="!(isNewItemModalActive || isEditItemModalActive)"
          :func="ActivateNewItemModal"
        />
        <div class="mt-6">
          <div>
            <h1 class="title is-3">Items</h1>
            <div v-if="viewItemsDisplay.length > 0">
              <label class="label">Search for items</label>
              <b-field class="mb-3">
                <b-input
                  icon="magnify"
                  size="is-medium"
                  placeholder="Enter a item name"
                  type="search"
                  expanded
                  v-model="itemSearch"
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
            </div>
            <b-loading
              :is-full-page="false"
              :active.sync="pageLoading"
              :can-cancel="false"
            ></b-loading>
            <section>
              <div
                class="card pointer-cursor-on-hover mb-3"
                @click="ActivateNewItemModal"
              >
                <div class="card-content">
                  <div class="media">
                    <div class="media-left">
                      <b-icon icon="currency-usd" size="is-medium"> </b-icon>
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
            <b-table
              :data="viewItemsDisplay"
              :columns="viewTableColumns"
              bordered
              focusable
              striped
              paginated
              :checkable="itemEditing"
              :checked-rows.sync="itemCheckedRows"
              per-page="20"
              :selected.sync="selected"
              v-if="viewItemsDisplay.length > 0"
            >
            </b-table>
            <b-field grouped class="mt-3 mb-3">
              <b-button
                v-if="itemEditing === true"
                icon-left="delta"
                native-type="submit"
                :expanded="deviceIsMobile"
                :loading="submitLoading"
                :disabled="submitLoading"
                @click="itemEditing = false"
              >
                Cancel editing
              </b-button>
              <b-button
                v-if="itemEditing === false"
                icon-left="delta"
                native-type="submit"
                :expanded="deviceIsMobile"
                :loading="submitLoading"
                :disabled="submitLoading"
                @click="itemEditing = true"
              >
                Edit items
              </b-button>
              <b-button
                type="is-danger"
                icon-left="delete"
                native-type="submit"
                :loading="pageLoading"
                :disabled="itemCheckedRows.length < 1"
                @click="DeleteCostsItems"
                >Delete selection
              </b-button>
            </b-field>
          </div>
          <p>{{ viewItemsDisplay.length }} entries displayed</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import costs from '@/requests/authenticated/costs'
import flatmates from '@/requests/authenticated/flatmates'
import cani from '@/requests/authenticated/can-i'
import { DialogProgrammatic as Dialog } from 'buefy'
import { Chart } from 'chart.js/auto'

export default {
  name: 'costs-view',
  data () {
    return {
      viewTableColumns: [
        {
          field: 'title',
          label: 'Title'
        },
        {
          field: 'frequency',
          label: 'Frequency'
        },
        {
          field: 'amount',
          label: 'Amount'
        },
        {
          field: 'invoiceDate',
          label: 'Invoice Date'
        }
      ],
      displayFloatingAddButton: true,
      canUserAccountAdmin: false,
      notes: '',
      users: {},
      view: {},
      chart: null,
      chartData: {
        type: 'bar',
        data: {
          datasets: [
            {
              label: 'Spending over the months',
              data: []
            }
          ]
        },
        options: {
          responsive: true,
          lineTension: 0.4,
          scales: {
            y: {
              ticks: {
                callback: (value, index, ticks) => '$' + value,
                stepSize: 10
              }
            }
          }
        }
      },
      deviceIsMobile: false,
      pageLoading: true,
      sortBy: 'recentlyUpdated',
      isNewItemModalActive: false,
      isEditItemModalActive: false,
      editItemProps: {
        id: undefined
      },
      selected: null,
      itemEditing: false,
      itemCheckedRows: [],
      deleteLoading: false
    }
  },
  components: {
    floatingAddButton: () =>
      import('@/components/common/floating-add-button.vue'),
    costsItemNew: () =>
      import('@/components/authenticated/costs-item-new.vue'),
    costsItemEdit: () =>
      import('@/components/authenticated/costs-item-edit.vue')
  },
  computed: {
    viewItemsDisplay () {
      // TODO express invoicedDate, invoicedBy, author, authorLast
      // TODO use this.users as id-to-names cache
      return (
        this.viewItemsFiltered.map((i) => {
          if (i.invoiceDate !== 0) {
            i.invoiceDate = common.TimestampToCalendarDate(i.invoiceDate)
          }
          i.amount = common.FormatFloatAsMoney(i.amount)
          return i
        }) || []
      )
    },
    viewItemsFiltered () {
      return this.view.items || []
    }
  },
  methods: {
    ActivateNewItemModal () {
      this.isNewItemModalActive = true
    },
    ActivateEditItemModal (itemId) {
      this.editItemProps = {
        id: itemId
      }
      this.isEditItemModalActive = true
    },
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    GetView () {
      costs
        .GetView()
        .then((resp) => {
          this.pageLoading = false
          this.view = resp.data.data || {}
        })
        .catch(() => {
          common.DisplayFailureToast(
            'Hmmm seems somethings gone wrong loading the costs view'
          )
        })
    },
    GetCostsNotes () {
      costs
        .GetCostsNotes()
        .then((resp) => {
          this.notes = resp.data.spec.notes || ''
          this.pageLoading = false
        })
        .catch(() => {
          common.DisplayFailureToast(
            'Hmmm seems somethings gone wrong loading the notes for costs'
          )
        })
    },
    EditCostsNotes () {
      if (this.canUserAccountAdmin !== true) {
        return
      }
      this.displayFloatingAddButton = false
      Dialog.prompt({
        title: 'Cost notes',
        message: `Enter notes that are useful for keeping track of costs.`,
        container: null,
        icon: 'text',
        hasIcon: true,
        inputAttrs: {
          placeholder: 'e.g. make sure to add rent costs.',
          maxlength: 250,
          required: false,
          value: this.notes || undefined
        },
        trapFocus: true,
        onConfirm: (value) => {
          costs
            .PutCostsNotes(value)
            .then(() => {
              this.pageLoading = true
              this.displayFloatingAddButton = true
              this.GetCostsNotes()
            })
            .catch((err) => {
              common.DisplayFailureToast(
                'Failed to update notes' +
                    `<br/>${err.response.data.metadata.response}`
              )
              this.displayFloatingAddButton = true
            })
        },
        onCancel: () => {
          this.displayFloatingAddButton = true
        }
      })
    },
    GetFlatmateName (id) {
      flatmates
        .GetFlatmate(id)
        .then((resp) => {
          return resp.data.spec.names
        })
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to fetch user account' +
                `<br/>${err.response.data.metadata.response}`
          )
          return id
        })
    },
    FmtAsCurrency (input) {
      const fixed = input.toFixed()
      const result = new Intl.NumberFormat().format(fixed)
      return result
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    },
    buildChart () {
      // TODO handle null values: e.g new instance
      if (
        this.view === null ||
          this.view.totalPriceMonthlyCumulative === null
      ) {
        return
      }
      const ctx = document.getElementById('costs-view-chart')
      this.chartData.data.datasets[0].data =
          this.view.totalPriceMonthlyCumulative.map((i) => i.total)
      this.chartData.data.labels = this.view.totalPriceMonthlyCumulative.map(
        (i) => i.yearMonth
      )
      if (this.chart !== null) {
        this.chart.destroy()
      }
      this.chart = new Chart(ctx, this.chartData)
    },
    DeleteCostsItems () {
      if (this.itemCheckedRows.length < 1) {
        return
      }
      const ids = this.itemCheckedRows.map((i) => i.id)
      Dialog.confirm({
        title: `Delete ${this.itemCheckedRows.length} item(s)`,
        message:
            'Are you sure that you wish to delete these cost items?' +
            '<br/>' +
            'This action cannot be undone.',
        confirmText: 'Delete item',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true
          costs
            .DeleteCostsItems(ids)
            .then((resp) => {
              common.DisplaySuccessToast(resp.data.metadata.response)
              this.GetView()
            })
            .catch((err) => {
              this.deleteLoading = false
              common.DisplayFailureToast(
                'Failed to delete cost items' +
                    ' - ' +
                    err.response.data.metadata.response
              )
            })
        }
      })
    }
  },
  watch: {
    selected () {
      if (this.selected === null) {
        return
      }
      // TODO launch edit modal
      this.ActivateEditItemModal(this.selected.id)
      this.selected = null
    },
    isNewItemModalActive () {
      if (this.isNewItemModalActive === false) {
        this.GetView()
      }
    },
    isEditItemModalActive () {
      if (this.isEditItemModalActive === false) {
        this.GetView()
      }
    },
    view () {
      if (this.view === null) {
        return
      }
      this.buildChart()
    }
  },
  async beforeMount () {
    cani.GetCanIgroup('admin').then((resp) => {
      this.canUserAccountAdmin = resp.data.data
    })
    this.GetView()
    this.GetCostsNotes()
  },
  mounted () {},
  beforeDestroy () {
    window.removeEventListener('resize', this.CheckDeviceIsMobile, true)
  },
  async created () {
    this.CheckDeviceIsMobile()
    window.addEventListener('resize', this.CheckDeviceIsMobile, true)
  }
}
</script>

<style scoped></style>
