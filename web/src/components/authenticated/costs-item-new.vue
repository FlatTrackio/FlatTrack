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
  <div class="item-page">
    <div class="modal-card" style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">New cost</p>
        <p class="modal-card-subtitle">Add a cost to track</p>
      </header>
      <section class="modal-card-body">
        <div>
          <b-field label="Name" class="is-marginless">
            <b-input
              type="text"
              v-model="title"
              size="is-medium"
              maxlength="30"
              icon="text"
              placeholder="Enter a title for the cost item"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="title = ''"
              @keyup.enter.native="PostCostsItem"
              autofocus
              required
            >
            </b-input>
          </b-field>
          <b-field label="Notes (optional)" class="is-marginless">
            <b-input
              type="text"
              v-model="notes"
              size="is-medium"
              icon="text"
              placeholder="Enter information extra"
              @keyup.enter.native="PostCostsItem"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="notes = ''"
              maxlength="40"
            >
            </b-input>
          </b-field>
          <b-field label="Amount">
            <b-input
              type="number"
              step="0.1"
              placeholder="0.00"
              v-model="amount"
              icon="currency-usd"
              icon-right="close-circle"
              icon-right-clickable
              required
              expanded
              @icon-right-click="amount = undefined"
              @keyup.enter.native="PostCostsItem"
              size="is-medium"
            >
            </b-input>
          </b-field>
          <b-field label="Frequency">
            <b-select
              v-model="frequency"
              placeholder="Select payment type"
              size="is-medium"
              icon="calendar"
              expanded
            >
              <option value="">Never</option>
              <option disabled value="reoccuring">Reoccurance</option>
              <option value="daily">Daily from invoice date</option>
              <option value="weekly">Weekly from invoice date</option>
              <option value="fortnightly">Fortnightly from invoice date</option>
              <option value="monthly">Monthly from invoice date</option>
            </b-select>
          </b-field>
          <b-field label="Invoice Date">
            <b-datepicker v-model="jsInvoiceDate" inline required>
            </b-datepicker>
          </b-field>
          <b-field
            label="Reoccur cost until (optional)"
            grouped
            group-multiline
            v-if="frequency !== ''"
          >
            <b-datepicker v-model="jsReoccurUntil" inline required>
            </b-datepicker>
            <b-button type="text" @click="jsReoccurUntil = null"
              >Clear</b-button
            >
          </b-field>
          <b-field label="Invoice link (optional)" class="is-marginless">
            <b-input
              type="text"
              v-model="invoiceLink"
              size="is-medium"
              icon="link"
              placeholder="Enter a link to download the invoice"
              @keyup.enter.native="PostCostsItem"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="invoiceLink = ''"
              maxlength="100"
            >
            </b-input>
          </b-field>
          <div
            class="columns"
            v-if="title !== '' && amount !== undefined && jsInvoiceDate !== undefined && frequency !== ''"
          >
            <div class="column">
              <p class="text">Will next occur: {{ NextOccurance }}</p>
              <p
                v-if="frequency !== '' && jsInvoiceDate !== undefined && (new Date(jsInvoiceDate)).getTime() < (new Date()).getTime()"
              >
                <b>Please note</b>: Reoccuring costs with past invoice dates
                will create subsequent cost items.
              </p>
            </div>
          </div>
          <b-field grouped>
            <b-button
              type="is-warning"
              size="is-medium"
              icon-left="arrow-left"
              native-type="submit"
              @click="$parent.close()"
            >
              Back
            </b-button>
            <b-button
              type="is-success"
              size="is-medium"
              icon-left="plus"
              native-type="submit"
              expanded
              :loading="submitLoading"
              :disabled="submitLoading"
              @click="PostCostsItem"
            >
              Add item
            </b-button>
          </b-field>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import costs from '@/requests/authenticated/costs'
import dayjs from 'dayjs'

export default {
  name: 'cost-item-new',
  data () {
    return {
      submitLoading: false,
      title: '',
      frequency: '',
      reoccurUntil: undefined,
      jsReoccurUntil: undefined,
      notes: '',
      amount: undefined,
      invoiceDate: undefined,
      invoiceLink: undefined,
      jsInvoiceDate: undefined
    }
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    PostCostsItem () {
      this.submitLoading = true
      if (this.notes === '') {
        this.notes = undefined
      }
      if (this.amount === 0) {
        this.amount = undefined
      } else {
        this.amount = parseFloat(this.amount)
      }

      if (this.jsInvoiceDate === 0) {
        common.DisplayFailureToast('Invoice date not set')
        this.submitLoading = false
        return
      }

      this.reoccurUntil = new Date(this.jsReoccurUntil || 0).getTime() / 1000
      this.invoiceDate = new Date(this.jsInvoiceDate || 0).getTime() / 1000

      if (
        this.invoiceLink !== '' &&
          typeof this.invoiceLink !== 'undefined' &&
          this.invoiceLink !== undefined
      ) {
        try {
          URL(this.invoiceLink)
        } catch {
          common.DisplayFailureToast('Invoice link invalid')
          this.submitLoading = false
          return
        }
      }

      costs
        .PostCostsItem({
          title: this.title,
          notes: this.notes,
          frequency: this.frequency,
          reoccurUntil: this.reoccurUntil,
          amount: this.amount,
          invoiceDate: this.invoiceDate,
          invoiceLink: this.invoiceLink
        })
        .then((resp) => {
          var item = resp.data.spec
          if (item.id !== '' || typeof item.id === 'undefined') {
            this.$parent.close()
          } else {
            this.submitLoading = false
            common.DisplayFailureToast('Unable to find created costs item')
          }
        })
        .catch((err) => {
          this.submitLoading = false
          common.DisplayFailureToast(
            `Failed to add costs item - ${err.response.data.metadata.response}`
          )
        })
    }
  },
  computed: {
    NextOccurance () {
      let nextOccurance
      let invoiceTimestamp = dayjs(this.jsInvoiceDate)
      switch (this.frequency) {
        case 'daily':
          nextOccurance = invoiceTimestamp.add(1, 'day')
          break
        case 'weekly':
          nextOccurance = invoiceTimestamp.add(1, 'week')
          break
        case 'fortnightly':
          nextOccurance = invoiceTimestamp.add(2, 'week')
          break
        case 'monthly':
          nextOccurance = invoiceTimestamp.add(1, 'month')
          break
        default:
          return
      }
      const calDate = common.TimestampToCalendarDate(nextOccurance.unix())
      return calDate
    }
  },
  async beforeMount () {}
}
</script>
