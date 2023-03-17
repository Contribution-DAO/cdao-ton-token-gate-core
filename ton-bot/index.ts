import { Telegraf } from "telegraf"
import { v4 as uuidV4 } from 'uuid'
import { config as dotenvConfig } from 'dotenv'
import axios from "axios"

dotenvConfig()

const bot = new Telegraf(process.env.BOT_TOKEN as string)

bot.start((ctx) => {
  let message = `Hello world`
  ctx.reply(message)
})

bot.on('chat_join_request', async (ctx) => {
  console.log("==============> new member joined")
  const chatId = ctx.chat?.id
  if (chatId) {
    console.log(process.env.API_HOST + "/internal/" + ctx.chatJoinRequest.user_chat_id + "/telegram/groups/" + chatId)
    const membership = await axios.get(process.env.API_HOST + "/internal/" + ctx.chatJoinRequest.user_chat_id + "/telegram/groups/" + chatId, {
      headers: {
        Authorization: "Bearer " + process.env.API_SECRET,
      }
    })

    if (membership.data.isMinted) {
      await ctx.approveChatJoinRequest(ctx.chatJoinRequest.user_chat_id)
      await axios.post(process.env.API_HOST + "/internal/" + ctx.chatJoinRequest.user_chat_id + "/telegram/groups/" + chatId + "/mark_join", {
        headers: {
          Authorization: "Bearer " + process.env.API_SECRET,
        }
      })
    } else {
      return await ctx.telegram.sendMessage(chatId, "Rejected " + ctx.chatJoinRequest.from.first_name + "! Please tell him to join again after minting his SBT.")
    }
    // return await ctx.telegram.sendMessage(chatId, "Pending " + ctx.chatJoinRequest.from.first_name + " (ID: " + ctx.chatJoinRequest.user_chat_id + ")")
  }
})

bot.on('new_chat_members', async (ctx) => {
  console.log("==============> new member added")
  const newMembers = ctx.message.new_chat_members
  const chatId = ctx.chat?.id
  const botId = (await bot.telegram.getMe()).id

  if (chatId) {
    if (newMembers.some(member => member.id === botId)) {
      // Bot is a new member of the chat
      bot.telegram.sendMessage(chatId, `Please input this group ID in Ton connect UI: ${chatId}`)
    } else {
      // New members have joined the chat
      newMembers.forEach(member => {
        bot.telegram.sendMessage(chatId, "Hello " + member.first_name + " (ID: " + member.id + ")")
      })
    }

    // return await ctx.telegram.sendMessage(chatId, "Hello " + newMembers[0].first_name + " (ID: " + newMembers[0].id + ")")
  }
});


bot.launch()

console.log("Bot launched")