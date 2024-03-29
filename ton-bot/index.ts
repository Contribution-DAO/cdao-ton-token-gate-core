import { Telegraf } from "telegraf"
import { v4 as uuidV4 } from 'uuid'
import { config as dotenvConfig } from 'dotenv'
import axios from "axios"
import TelegramBot from "node-telegram-bot-api"

dotenvConfig()

const bot = new TelegramBot(process.env.BOT_TOKEN as string, {polling: true});

// bot.start((ctx) => {
//   let message = `Hello world`
//   ctx.reply(message)
// })

bot.on('chat_join_request', async (ctx) => {
  let approved = false

  console.log("==============> new member joined")
  const chatId = ctx.chat?.id
  const userId = ctx.user_chat_id

  try {
    if (chatId) {
      console.log(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId)
      const membership = await axios.get(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId, {
        headers: {
          Authorization: "Bearer " + process.env.API_SECRET,
        }
      })
  
      if (membership.data.isMinted && !membership.data.isJoined) {
        await bot.approveChatJoinRequest(chatId, userId)
        approved = true
        await axios.post(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId + "/mark_joined", {
          joined: true,
        }, {
          headers: {
            Authorization: "Bearer " + process.env.API_SECRET,
          }
        })
      } else {
        await bot.declineChatJoinRequest(chatId, userId)
        // return await ctx.telegram.sendMessage(chatId, "Rejected " + ctx.chatJoinRequest.from.first_name + "! Please tell him to join again after minting his SBT.")
      }
      // return await ctx.telegram.sendMessage(chatId, "Pending " + ctx.chatJoinRequest.from.first_name + " (ID: " + ctx.chatJoinRequest.user_chat_id + ")")
    }
  } catch (err) {
    if (!approved) {
      try {
        await bot.declineChatJoinRequest(chatId, userId)
      } catch (err) {
        console.error(err)
      }
    }
  }
})

// bot.on('chat_member', async (ctx) => {
//   let approved = false;

//   console.log("==============> new member updated")
//   const chatId = ctx.chat?.id
//   const userId = ctx.chatMember.new_chat_member.user.id

//   try {
//     if (chatId) {
//       console.log(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId)
//       const membership = await axios.get(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId, {
//         headers: {
//           Authorization: "Bearer " + process.env.API_SECRET,
//         }
//       })
  
//       if (membership.data.isMinted && !membership.data.isJoined) {
//         await ctx.approveChatJoinRequest(userId)
//         approved = true
//         await axios.post(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId + "/mark_joined", {
//           joined: true,
//         }, {
//           headers: {
//             Authorization: "Bearer " + process.env.API_SECRET,
//           }
//         })
//       } else {
//         await ctx.declineChatJoinRequest(userId)
//         // return await ctx.telegram.sendMessage(chatId, "Rejected " + ctx.chatJoinRequest.from.first_name + "! Please tell him to join again after minting his SBT.")
//       }
//       // return await ctx.telegram.sendMessage(chatId, "Pending " + ctx.chatJoinRequest.from.first_name + " (ID: " + ctx.chatJoinRequest.user_chat_id + ")")
//     }
//   } catch (err) {
//     if (!approved) {
//       await ctx.declineChatJoinRequest(userId)
//     }
//   }
// })

bot.on('left_chat_member', async (ctx) => {
  try {
    console.log("==============> member left")
    const member = ctx.left_chat_member
    const chatId = ctx.chat?.id

    if (!member) return;
  
    await axios.post(process.env.API_HOST + "/internal/" + member.id + "/telegram/groups/" + chatId + "/mark_joined", {
      joined: false,
    }, {
      headers: {
        Authorization: "Bearer " + process.env.API_SECRET,
      }
    })
  } catch (err) {

  }

})

bot.on('new_chat_members', async (ctx) => {
  let approved = false

  console.log("==============> new member added")
  const newMembers = ctx.new_chat_members ?? []
  const chatId = ctx.chat?.id
  const botId = (await bot.getMe()).id
  
  if (chatId) {
    for (const member of newMembers) {
      try {
        if (member.id === botId) {
          // Bot is a new member of the chat
          await bot.sendMessage(chatId, `Please input this group ID in Ton connect UI: ${chatId}`)
        } else {
          // New members have joined the chat
          const membership = await axios.get(process.env.API_HOST + "/internal/" + member.id + "/telegram/groups/" + chatId, {
            headers: {
              Authorization: "Bearer " + process.env.API_SECRET,
            }
          })

          if (membership.data.isMinted) {
            // await bot.telegram.sendMessage(chatId, "Hello " + member.first_name /*+ " (ID: " + member.id + ")"*/)

            approved = true
            await axios.post(process.env.API_HOST + "/internal/" + member.id + "/telegram/groups/" + chatId + "/mark_joined", {
              joined: true,
            }, {
              headers: {
                Authorization: "Bearer " + process.env.API_SECRET,
              }
            })
          } else {
            await bot.banChatMember(chatId, member.id)
          }
        }
      } catch (err) {
        if (!approved) {
          try {
            await bot.banChatMember(chatId, member.id)
          } catch (err) {
            console.error(err)
          }
        }
      }
    }

    // return await ctx.telegram.sendMessage(chatId, "Hello " + newMembers[0].first_name + " (ID: " + newMembers[0].id + ")")
  }

});

console.log("Bot launched")