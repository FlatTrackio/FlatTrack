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
    <section>
      <div class="card pointer-cursor-on-hover">
        <div class="card-content">
          <div class="media">
            <div
              class="media-left"
              @click="
                $router.push({
                  name: 'Edit task',
                  params: { id: task.id },
                })
              "
            >
              <b-icon
                icon="clipboard-list"
                :type="task.completed === true ? 'is-success' : ''"
                size="is-medium"
              >
              </b-icon>
            </div>
            <div
              class="media-content"
              @click="
                $router.push({
                  name: 'Edit task',
                  params: { id: task.id },
                })
              "
            >
              <div class="display-items-on-the-same-line">
                <p class="title is-4">{{ task.name }}</p>
              </div>
              <p class="subtitle is-6">
                <span
                  v-if="task.creationTimestamp == task.modificationTimestamp"
                >
                  Created {{ TimestampToCalendar(task.creationTimestamp) }}
                </span>
                <span v-else>
                  Updated {{ TimestampToCalendar(task.modificationTimestamp) }}
                </span>
              </p>
            </div>
            <div class="media-right">
              <b-tooltip label="Delete" class="is-paddingless" :delay="200">
                <b-button
                  type="is-danger"
                  icon-right="delete"
                  :loading="itemDeleting"
                  v-if="deviceIsMobile === false"
                  @click="DeleteShoppingList(task.id)"
                />
              </b-tooltip>
              <b-icon
                icon="chevron-right"
                size="is-medium"
                type="is-midgray"
              ></b-icon>
            </div>
          </div>
          <div
            v-if="!mini"
            class="content"
            @click="
              $router.push({
                name: 'Edit task',
                params: { id: task.id },
              })
            "
          >
            <div>
              <b-tag :type="task.completed ? 'is-info' : 'is-warning'">
                {{ task.completed ? "Completed" : "Uncompleted" }}
              </b-tag>
            </div>
            <br />
            <span v-if="task.notes !== '' && typeof task.notes !== 'undefined'">
              <i>
                {{ PreviewNotes(task.notes) }}
              </i>
              <br />
              <br />
            </span>
          </div>
        </div>
      </div>
      <br />
    </section>
  </div>
</template>

<script>
import { DialogProgrammatic as Dialog } from "buefy";
import common from "@/common/common";
import shoppinglist from "@/requests/authenticated/shoppinglist";
import shoppinglistCommon from "@/common/shoppinglist";

export default {
  name: "task-card-view",
  props: {
    mini: Boolean,
    deviceIsMobile: Boolean,
    index: Number,
    tasks: Object,
    task: Object,
  },
  data() {
    return {
      deleteLoading: false,
      authorNames: "",
      authorLastNames: "",
    };
  },
  methods: {
    PreviewNotes(notes) {
      if (notes.length <= 35) {
        return notes;
      }
      var notesBytes = notes.split("");
      var notesBytesValid = notesBytes.filter((value, index) => {
        if (index <= 35) {
          return value;
        }
      });
      return notesBytesValid.join("") + "...";
    },
    DeleteShoppingList(id) {
      Dialog.confirm({
        title: "Delete shopping list",
        message:
          "Are you sure that you wish to delete this shopping list?" +
          "<br/>" +
          "This action cannot be undone.",
        confirmText: "Delete shopping list",
        type: "is-danger",
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true;
          window.clearInterval(this.intervalLoop);
          shoppinglist
            .DeleteShoppingList(id)
            .then((resp) => {
              let removedFromLists = this.lists;
              removedFromLists.splice(this.index, 1);
              this.$emit("lists", removedFromLists);

              common.DisplaySuccessToast("Deleted the shopping list");
              shoppinglistCommon.DeleteShoppingListFromCache(id);
            })
            .catch((err) => {
              this.deleteLoading = false;
              common.DisplayFailureToast(
                "Failed to delete the shopping list" +
                  "<br/>" +
                  err.response.data.metadata.response
              );
            });
        },
      });
    },
    TimestampToCalendar(timestamp) {
      return common.TimestampToCalendar(timestamp);
    },
  },
  async beforeMount() {},
};
</script>

<style>
.display-items-on-the-same-line {
  display: flex;
}

.display-items-on-the-same-line div {
  margin-left: 10px;
}
</style>
